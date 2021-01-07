package service

import "log"

func (h *HermesService) StartListeningForMessages() error {
	log.Printf("listening for pubSub Messges...")
	err := h.pubSub.PullMsgs()
	if err != nil {
		return err
	}
	return nil
}