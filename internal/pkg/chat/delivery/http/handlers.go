package http

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"yula/internal/models"
	"yula/internal/pkg/advt"
	"yula/internal/pkg/logging"
	"yula/internal/pkg/middleware"
	"yula/internal/pkg/user"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"

	proto "yula/proto/generated/chat"

	internalError "yula/internal/error"
)

var (
	logger       = logging.GetLogger()
	chatSessions = map[string]*ChatSession{} // to_string(idFrom) + "->" + to_string(idTo) + ":" + to_string(idAdv) => conn
)

type ChatSession struct {
	idFrom int64
	idTo   int64
	idAdv  int64

	conn []*websocket.Conn
}

type ChatHandler struct {
	cu proto.ChatClient
	au advt.AdvtUsecase
	uu user.UserUsecase
}

func NewChatHandler(cu proto.ChatClient, au advt.AdvtUsecase, uu user.UserUsecase) *ChatHandler {
	return &ChatHandler{
		cu: cu,
		au: au,
		uu: uu,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (ch *ChatHandler) Routing(r *mux.Router, sm *middleware.SessionMiddleware) {
	s := r.PathPrefix("/chat").Subrouter()
	s.HandleFunc("/connect/{idFrom:[0-9]+}/{idTo:[0-9]+}/{idAdv:[0-9]+}", middleware.SetSCRFToken(http.HandlerFunc(ch.ConnectHandler))).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/createDialog/{idFrom:[0-9]+}/{idTo:[0-9]+}/{idAdv:[0-9]+}", middleware.SetSCRFToken(http.HandlerFunc(ch.CreateDialog))).Methods(http.MethodPost, http.MethodOptions)

	s.HandleFunc("/getDialogs/{idFrom:[0-9]+}", middleware.SetSCRFToken(sm.CheckAuthorized(http.HandlerFunc(ch.getDialogsHandler)))).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/getHistory/{idFrom:[0-9]+}/{idTo:[0-9]+}/{idAdv:[0-9]+}", middleware.SetSCRFToken(sm.CheckAuthorized(http.HandlerFunc(ch.getHistoryHandler)))).Methods(http.MethodGet, http.MethodOptions)

	s.Handle("/clear/{idFrom:[0-9]+}/{idTo:[0-9]+}/{idAdv:[0-9]+}", sm.CheckAuthorized(http.HandlerFunc(ch.ClearHandler))).Methods(http.MethodPost, http.MethodOptions)
}

func (ch *ChatHandler) CreateDialog(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))

	vars := mux.Vars(r)
	idFrom, _ := strconv.ParseInt(vars["idFrom"], 10, 64)
	idTo, err := strconv.ParseInt(vars["idTo"], 10, 64)
	idAdv, _ := strconv.ParseInt(vars["idAdv"], 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	_, err = ch.cu.CreateDialog(context.Background(), &proto.Dialog{
		DI: &proto.DialogIdentifier{
			Id1:   idFrom,
			Id2:   idTo,
			IdAdv: idAdv,
		},
		CreatedAt: timestamppb.Now(),
	})

	if err != nil {
		logger.Warnf("create dialog error: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(models.ToBytes(http.StatusOK, "dialog create success", nil))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
	logger.Info("dialog create success")
}

func (ch *ChatHandler) ConnectHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))

	vars := mux.Vars(r)
	idFrom, _ := strconv.ParseInt(vars["idFrom"], 10, 64)
	idTo, _ := strconv.ParseInt(vars["idTo"], 10, 64)
	idAdv, _ := strconv.ParseInt(vars["idAdv"], 10, 64)

	websocketConnection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Info("can not upgrade connection to websocket: ", err.Error())
		return
	}

	var curSession *ChatSession
	key := fmt.Sprintf("%d->%d:%d", idFrom, idTo, idAdv)
	if val, ok := chatSessions[key]; ok {
		curSession = val
		val.conn = append(val.conn, websocketConnection)
	} else {
		curSession = &ChatSession{
			idFrom: idFrom,
			idTo:   idTo,
			idAdv:  idAdv,

			conn: []*websocket.Conn{websocketConnection},
		}
		chatSessions[key] = curSession
	}

	go ch.HandleMessages(curSession, websocketConnection)
}

func (ch *ChatHandler) HandleMessages(session *ChatSession, conn *websocket.Conn) {
	defer func() {
		conn.Close()

		key := fmt.Sprintf("%d->%d:%d", session.idFrom, session.idTo, session.idAdv)
		for ind, value := range chatSessions[key].conn {
			if value == conn {
				chatSessions[key].conn[ind] = chatSessions[key].conn[len(chatSessions[key].conn)-1]
				chatSessions[key].conn[len(chatSessions[key].conn)-1] = nil
				chatSessions[key].conn = chatSessions[key].conn[:len(chatSessions[key].conn)-1]
			}
		}
	}()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Debug(err)
			return
		}

		_, err = ch.cu.Create(context.Background(), &proto.Message{
			MI: &proto.MessageIdentifier{
				IdFrom: session.idFrom,
				IdTo:   session.idTo,
				IdAdv:  session.idAdv,
			},
			Msg:       string(msg),
			CreatedAt: timestamppb.Now(),
		})
		if err != nil {
			logger.Warnf("cannot create proto message %s", err.Error())
		}

		key := fmt.Sprintf("%d->%d:%d", session.idTo, session.idFrom, session.idAdv)
		to := chatSessions[key]

		if to == nil {
			continue
		}

		for _, conn := range to.conn {
			if err := conn.WriteMessage(msgType, msg); err != nil {
				logger.Errorf("Can not write msg from user %d to user %d on ad %d", session.idFrom, session.idTo, session.idAdv)
				return
			}
		}
	}
}

func (ch *ChatHandler) getHistoryHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))

	u, err := url.Parse(r.URL.RequestURI())
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(internalError.BadRequest)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	query := u.Query()
	page, _ := models.NewPage(query.Get("page"), query.Get("count"))

	vars := mux.Vars(r)
	idFrom, _ := strconv.ParseInt(vars["idFrom"], 10, 64)
	idTo, err := strconv.ParseInt(vars["idTo"], 10, 64)
	idAdv, _ := strconv.ParseInt(vars["idAdv"], 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	protomessages, err := ch.cu.GetHistory(context.Background(), &proto.GetHistoryArg{
		DI: &proto.DialogIdentifier{
			Id1:   idFrom,
			Id2:   idTo,
			IdAdv: idAdv,
		},
		FP: &proto.FilterParams{
			Offset: page.PageNum * page.Count,
			Limit:  page.Count,
		},
	})

	if err != nil {
		logger.Warnf("get history chat error: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.NotExist)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	var messages []*models.Message
	for _, message := range protomessages.M {
		messages = append(messages, &models.Message{
			MI: models.IMessage{
				IdFrom: message.MI.IdFrom,
				IdTo:   message.MI.IdTo,
				IdAdv:  message.MI.IdAdv,
			},
			Msg:       message.Msg,
			CreatedAt: message.CreatedAt.AsTime(),
		})
	}

	if err != nil {
		logger.Warnf("get history chat error: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyChatHistory{Messages: messages}
	_, err = w.Write(models.ToBytes(http.StatusOK, "chat history found successfully", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
	logger.Info("chat history found successfully")
}

func (ch *ChatHandler) ClearHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))

	vars := mux.Vars(r)
	idFrom, _ := strconv.ParseInt(vars["idFrom"], 10, 64)
	idTo, err := strconv.ParseInt(vars["idTo"], 10, 64)
	idAdv, _ := strconv.ParseInt(vars["idAdv"], 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	_, err = ch.cu.Clear(context.Background(), &proto.DialogIdentifier{
		Id1:   idFrom,
		Id2:   idTo,
		IdAdv: idAdv,
	})

	if err != nil {
		logger.Warnf("clear chat error: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(models.ToBytes(http.StatusOK, "clear chat success", nil))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
	logger.Info("clear chat success")
}

func (ch *ChatHandler) getDialogsHandler(w http.ResponseWriter, r *http.Request) {
	logger = logger.GetLoggerWithFields((r.Context().Value(middleware.ContextLoggerField)).(logrus.Fields))

	vars := mux.Vars(r)
	idFrom, err := strconv.ParseInt(vars["idFrom"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		metaCode, metaMessage := internalError.ToMetaStatus(err)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	protodialogs, err := ch.cu.GetDialogs(context.Background(), &proto.UserIdentifier{IdFrom: idFrom})
	if err != nil {
		logger.Warnf("get dialogs error: %s", err.Error())
		w.WriteHeader(http.StatusOK)

		metaCode, metaMessage := internalError.ToMetaStatus(internalError.NotExist)
		_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
		if err != nil {
			logger.Warnf("cannot write answer to body %s", err.Error())
		}
		return
	}

	var dialogs []*models.HttpDialog
	for _, dialog := range protodialogs.D {
		shortAd := &models.AdvertShort{
			Id:       -1,
			Name:     "dummy",
			Price:    -1,
			Location: "dummy",
			Image:    "dummy",
		}

		if dialog.DI.IdAdv != -1 {
			ad, err := ch.au.GetAdvert(dialog.DI.IdAdv, -1, false)
			if err != nil {
				logger.Warnf("get adv error: %s", err.Error())
				w.WriteHeader(http.StatusOK)

				metaCode, metaMessage := internalError.ToMetaStatus(internalError.NotExist)
				_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
				if err != nil {
					logger.Warnf("cannot write answer to body %s", err.Error())
				}
				return
			}

			shortAd = ad.ToShort()
		}

		user2, err := ch.uu.GetById(dialog.DI.Id2)
		if err != nil {
			logger.Warnf("get user error: %s", err.Error())
			w.WriteHeader(http.StatusOK)

			metaCode, metaMessage := internalError.ToMetaStatus(internalError.NotExist)
			_, err = w.Write(models.ToBytes(metaCode, metaMessage, nil))
			if err != nil {
				logger.Warnf("cannot write answer to body %s", err.Error())
			}
			return
		}

		dialogs = append(dialogs, &models.HttpDialog{
			Id:        user2.Id,
			Name:      user2.Name,
			Surname:   user2.Surname,
			Adv:       *shortAd,
			CreatedAt: dialog.CreatedAt.AsTime(),
		})
	}

	w.WriteHeader(http.StatusOK)
	body := models.HttpBodyDialogs{Dialogs: dialogs}
	_, err = w.Write(models.ToBytes(http.StatusOK, "dialogs found successfully", body))
	if err != nil {
		logger.Warnf("cannot write answer to body %s", err.Error())
	}
	logger.Info("dialogs found successfully")
}
