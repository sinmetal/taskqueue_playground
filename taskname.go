package taskqueue_playground

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

func init() {
	api := TaskNameAPI{}

	http.HandleFunc("/", handler)
	http.HandleFunc("/taskname/withtx", api.withTx)
	http.HandleFunc("/taskname/withouttx", api.withoutTx)
}

type TaskNameAPI struct{}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func (api *TaskNameAPI) withoutTx(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	task := taskqueue.Task{
		Name:   "withouttx",
		Path:   "/",
		Method: http.MethodGet,
	}
	_, err := taskqueue.Add(c, &task, "")
	if err != nil {
		log.Errorf(c, "Taskqueue.Add failed: %v", err)
		http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("done."))
}

func (api *TaskNameAPI) withTx(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	err := datastore.RunInTransaction(c, func(c context.Context) error {
		task := taskqueue.Task{
			Name:   "withtx",
			Path:   "/",
			Method: http.MethodGet,
		}
		_, err := taskqueue.Add(c, &task, "")
		if err != nil {
			return err
		}

		return nil
	}, nil)
	if err != nil {
		log.Errorf(c, "Transaction failed: %v", err)
		http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("done."))
}
