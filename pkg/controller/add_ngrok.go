package controller

import (
	"github.com/zufardhiyaulhaq/ngrok-operator/pkg/controller/ngrok"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, ngrok.Add)
}
