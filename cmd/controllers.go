package cmd

import (
	"time"

	"github.com/dmathieu/dice/cloudprovider"
	"github.com/dmathieu/dice/controllers"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

type watchControllers struct {
	evictDoneCh     chan struct{}
	evictFinishedCh chan struct{}

	deleteDoneCh   chan struct{}
	informerDoneCh chan struct{}
}

func (w *watchControllers) Run(maxUptime *time.Duration) {
	start := time.Now()
	go func() {
		if maxUptime == nil {
			return
		}
		for {
			if time.Since(start) > *maxUptime {
				close(w.evictFinishedCh)
				w.Close()
				return
			}
		}
	}()

	<-w.evictFinishedCh
}

func (w *watchControllers) Close() {
	close(w.evictDoneCh)
	close(w.deleteDoneCh)
	close(w.informerDoneCh)
}

func runWatchControllers(k8Client kubernetes.Interface, cloudClient cloudprovider.CloudProvider, c int, i bool) (*watchControllers, error) {
	w := &watchControllers{}

	inform := informers.NewSharedInformerFactory(k8Client, time.Second*30)
	evict := controllers.NewEvictNodeController(k8Client, inform.Core().V1().Nodes(), c, i)
	w.evictDoneCh = make(chan struct{})
	w.evictFinishedCh = make(chan struct{})
	go evict.Run(w.evictFinishedCh, w.evictDoneCh)

	delete := controllers.NewDeleteNodeController(k8Client, cloudClient, inform.Core().V1().Pods(), inform.Core().V1().Nodes())
	w.deleteDoneCh = make(chan struct{})
	go delete.Run(w.deleteDoneCh)

	w.informerDoneCh = make(chan struct{})
	inform.Start(w.informerDoneCh)

	return w, nil
}
