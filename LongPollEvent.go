package VkApi

import (
	"errors"
	"strconv"
)

const LPMessageReplaceFlag = 1
const LPMessageSetFlag = 2
const LPMessageDropFlag = 3
const LPNewMessage = 4
const LPEditMessage = 5
const LPReadInputMessage = 6
const LPReadOutputMessage = 7
const LPOnline = 8
const LPOffline = 9
const LPDialogDropFlag = 10
const LPDialogReplaceFlang = 11
const LPDialogSetFlag = 12
const LPDeleteMessage = 13
const LPRestoreMessage = 14
const LPChangeGroupDialogSettings = 51
const LPTyping = 61
const LPTypingAtGroupDialog = 62
const LPCall = 70
const LPCounter = 80
const LPNotifications = 114

const FLAG_OUTBOX = 2
const FLAG_UNREAD = 1

type ActorAndUser struct {
	Actor int
	User  int
}

type MessageEvent struct {
	MessageId       int
	Mask            int
	PeerId          int
	Timestamp       int
	Text            string
	RandomId        int
	ChatCreate      bool
	ChatTitleUpdate bool
	ChatPhotoUpdate bool
	ChatInviteUser  *ActorAndUser
	ChatKilUser     *ActorAndUser
	Title           string
	IsEditMessage   bool
	Emoji           bool
	FromAdmin       int
	Geo             string
	GeoProvider     string
	From            int
}

func (m *MessageEvent) Fill(data Update) error {
	if len(data) >= 6 {
		if d, ok := data[1].(float64); ok {
			m.MessageId = int(d)
		} else {
			return errors.New("No message id\n")
		}

		if d, ok := data[2].(float64); ok {
			m.Mask = int(d)
		} else {
			return errors.New("No message mask\n")
		}

		if d, ok := data[3].(float64); ok {
			m.PeerId = int(d)
		} else {
			return errors.New("No message peer id \n")
		}

		if d, ok := data[4].(float64); ok {
			m.Timestamp = int(d)
		} else {
			return errors.New("No message timestamp\n")
		}

		if d, ok := data[5].(string); ok {
			m.Text = d
		} else {
			return errors.New("No message text\n")
		}
		if len(data) >= 7 {
			if d, ok := data[6].(map[string]interface{}); ok {
				if err := m.fillAttachments(d); err != nil {
					return err
				}
			} else {
				return errors.New("Bad attachments type\n")
			}
		}
		if len(data) >= 9 {
			if d, ok := data[7].(float64); ok {
				m.RandomId = int(d)
			} else {
				return errors.New("No message random id\n")
			}
		}
		return nil
	} else {
		return errors.New("Bad message size\n")
	}
}

func (m *MessageEvent) fillAttachments(attachments map[string]interface{}) error {
	if len(attachments) > 0 {
		for key, value := range attachments {
			if v, ok := value.(string); ok && key == "source_act" {
				if v == "chat_create" {
					m.ChatCreate = true
				} else if v == "chat_title_update" {
					m.ChatTitleUpdate = true
				} else if v == "chat_photo_update" {
					m.ChatPhotoUpdate = true
				} else if v == "chat_invite_user" {
					user, hasUser := attachments["source_mid"]
					actor, hasActor := attachments["from"]
					if hasUser && hasActor {
						uId, u := user.(string)
						aId, a := actor.(string)
						if a && u {
							uIdInt, errU := strconv.Atoi(uId)
							aIdInt, errA := strconv.Atoi(aId)
							if errA == nil && errU == nil {
								m.ChatInviteUser = &ActorAndUser{int(aIdInt), int(uIdInt)}
							} else if errU != nil {
								return errU
							} else if errA != nil {
								return errA
							}
						}
					}
				} else if v == "chat_kick_user" {
					user, hasUser := attachments["source_mid"]
					actor, hasActor := attachments["from"]
					if hasUser && hasActor {
						uId, u := user.(string)
						aId, a := actor.(string)
						if a && u {
							uIdInt, errU := strconv.Atoi(uId)
							aIdInt, errA := strconv.Atoi(aId)
							if errA == nil && errU == nil {
								m.ChatKilUser = &ActorAndUser{int(aIdInt), int(uIdInt)}
							} else if errU != nil {
								return errU
							} else if errA != nil {
								return errA
							}
						}
					}
				}
			}
			if v, ok := value.(string); ok && key == "title" {
				m.Title = v
			}
			if _, ok := value.(float64); ok && key == "emoji" {
				m.Emoji = true
			}
			if userId, ok := value.(float64); ok && key == "from_admin" {
				m.FromAdmin = int(userId)
			}
			if geo, ok := value.(string); ok && key == "geo" {
				m.Geo = geo
			}
			if geoProvider, ok := value.(string); ok && key == "geo_provider" {
				m.GeoProvider = geoProvider
			}
			if from, ok := value.(string); ok && key == "from" {
				fromId, err := strconv.Atoi(from)
				if err != nil {
					return err
				}
				m.From = int(fromId)
			}
		}
	}
	return nil
}
func (m *MessageEvent) IsOutMessage() bool {
	return m.Mask&FLAG_OUTBOX == FLAG_OUTBOX
}
