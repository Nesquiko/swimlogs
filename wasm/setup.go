package main

import (
	"github.com/Nesquiko/swimlogs/pkg/view/pages"
	"github.com/Nesquiko/swimlogs/pkg/view/state"
	"github.com/vugu/vgrouter"
	"github.com/vugu/vugu"
)

func vuguSetup(buildEnv *vugu.BuildEnv, eventEnv vugu.EventEnv) vugu.Builder {

	tss := state.TrainingStateStorage{}

	router := vgrouter.New(eventEnv)

	// MAKE OUR WIRE FUNCTION POPULATE ANYTHING THAT WANTS A "NAVIGATOR".
	buildEnv.SetWireFunc(func(b vugu.Builder) {
		if c, ok := b.(vgrouter.NavigatorSetter); ok {
			c.NavigatorSet(router)
		}
		if c, ok := b.(state.TrainingStateStorageSetter); ok {
			c.TrainingStateStorageSet(&tss)
		}
	})

	// CREATE THE ROOT COMPONENT
	root := &pages.Root{}
	buildEnv.WireComponent(root) // WIRE IT
	router.MustAddRouteExact(
		"/",
		vgrouter.RouteHandlerFunc(func(*vgrouter.RouteMatch) { root.Body = &pages.LandingPage{} }),
	)

	router.MustAddRouteExact(
		"/training",
		vgrouter.RouteHandlerFunc(func(*vgrouter.RouteMatch) { root.Body = &pages.TrainingPage{} }),
	)

	router.MustAddRouteExact(
		"/add",
		vgrouter.RouteHandlerFunc(func(*vgrouter.RouteMatch) { root.Body = &pages.ActionsPage{} }),
	)

	router.MustAddRouteExact(
		"/add/session",
		vgrouter.RouteHandlerFunc(
			func(*vgrouter.RouteMatch) { root.Body = &pages.SessionAddPage{} },
		),
	)

	// router.SetNotFound(vgrouter.RouteHandlerFunc(
	// 	func(rm *vgrouter.RouteMatch) {
	// 		root.Body = &pages.PageNotFound{} // A PAGE FOR THE NOT-FOUND CASE
	// 	}))

	// TELL THE ROUTER TO LISTEN FOR THE BROWSER CHANGING URLS
	err := router.ListenForPopState()
	if err != nil {
		panic(err)
	}

	// GRAB THE CURRENT BROWSER URL AND PROCESS IT AS A ROUTE
	err = router.Pull()
	if err != nil {
		panic(err)
	}

	return root
}
