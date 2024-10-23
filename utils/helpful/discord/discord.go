package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type LogType int

const (
	StartLog LogType = iota
	StripeLog
	ErrorLog
	MessageLog
)

func (l LogType) String() string {
	switch l {
	case StartLog:
		return "[START]"
	case StripeLog:
		return "[STRIPE]"
	case ErrorLog:
		return "[ERROR]"
	case MessageLog:
		return "[MESSAGE]"
	default:
		return "[UNKNOWN]"
	}
}

func getWebhookURL(logType LogType) string {
	switch logType {
	case StartLog:
		return "https://discord.com/api/webhooks/1287514708177457225/yJAQvIkzkUfzg_jeIr-SVjSfCCHWHtpRTXfsvRYGEZBTrD5iKrOo51hmesvZv2Oa5UBp"
	case StripeLog:
		return "https://discord.com/api/webhooks/1289951157007028346/ELDvDZVjjGXq_boxyrRa9VIpo104ZWRIK8zmyKIt6pH8J9lCNZcQY3MZdO2DwYhHPTFC"
	case ErrorLog:
		return "https://discord.com/api/webhooks/1289951002115706932/iac6aw5gHTiKgyPP1wrxgTta8iG761V4ivUMg1khXrgZ-LfzpZRodBCQLzxLbavmlhlK"
	case MessageLog:
		return "https://discord.com/api/webhooks/1289951086660161576/9UL38h_OJuaXZXR3-l43OMRhtczTbcHAjl2YXr-eMFYkQ62jLZureCNoyNHwtWOm83gU"
	default:
		return ""
	}
}

func SendDiscordWebhook(logType LogType, message string) error {
	webhookURL := getWebhookURL(logType)
	if webhookURL == "" {
		return fmt.Errorf("invalid log type: %v", logType)
	}

	fullMessage := fmt.Sprintf("%s %s", logType.String(), message)

	payload := map[string]string{
		"content": fullMessage,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received non-204 status code: %d", resp.StatusCode)
	}

	return nil
}

func SendMessage(logType LogType, message string) error {
	return SendDiscordWebhook(logType, message)
}
