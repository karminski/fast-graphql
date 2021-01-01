// subscription-executor.go

package backend

import (
    "sync"
)


type SubscriptionExecutor struct {
	mutex sync.RWMutex
	subs  map[string]chan string
}


func NewSubscriptionExecutor() *SubscriptionExecutor {
	subExec      := &SubscriptionExecutor{}
	subExec.subs  = make(map[string]chan string)
	return subExec
}
// 
// 
// func (subExec *SubscriptionExecutor) Subscribe(topic string, ch chan string) {
// 	subExec.mutex.Lock()
// 	defer subExec.mutex.Unlock()
// 
// 	subExec.subs[topic] = append(subExec.subs[topic], ch)
// }
// 
// 
// func (subExec *SubscriptionExecutor) PublishMutation(topic string, msg string) {
// 	subExec.mutex.Lock()
// 	defer subExec.mutex.Unlock()
// 
// 	for _, ch := range subExec.subs[topic] {
// 		ch <- msg
// 	}
// }
// 
// 
// func (subExec *SubscriptionExecutor) RegistSubscription(topic string) error {
// 	ch := make(chan string, 1)
// 	subExec.Subscribe(topic, ch)
// 	return ch
// } 

