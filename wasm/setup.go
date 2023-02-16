package main

import (
	"github.com/Nesquiko/swimlogs/pkg/view/pages"
	"github.com/Nesquiko/swimlogs/pkg/view/state"
	"github.com/vugu/vgrouter"
	"github.com/vugu/vugu"
)

func vuguSetup(buildEnv *vugu.BuildEnv, eventEnv vugu.EventEnv) vugu.Builder {

	tss := state.TrainingStateStorage{}
	sss := state.SessionStateStorage{}

	router := vgrouter.New(eventEnv)

	// MAKE OUR WIRE FUNCTION POPULATE ANYTHING THAT WANTS A "NAVIGATOR".
	buildEnv.SetWireFunc(func(b vugu.Builder) {
		if c, ok := b.(vgrouter.NavigatorSetter); ok {
			c.NavigatorSet(router)
		}
		if c, ok := b.(state.TrainingStateStorageSetter); ok {
			c.TrainingStateStorageSet(&tss)
		}
		if c, ok := b.(state.SessionStateStorageSetter); ok {
			c.SessionStateStorageSet(&sss)
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
		"/create/session",
		vgrouter.RouteHandlerFunc(
			func(*vgrouter.RouteMatch) { root.Body = &pages.SessionCreatePage{} },
		),
	)

	router.MustAddRouteExact(
		"/sessions",
		vgrouter.RouteHandlerFunc(
			func(*vgrouter.RouteMatch) { root.Body = &pages.SessionsPage{} },
		),
	)

	router.MustAddRouteExact(
		"/edit/session",
		vgrouter.RouteHandlerFunc(
			func(*vgrouter.RouteMatch) { root.Body = &pages.SessionEditPage{} },
		),
	)

	router.MustAddRouteExact(
		"/create/training",
		vgrouter.RouteHandlerFunc(
			func(*vgrouter.RouteMatch) { root.Body = &pages.TrainingCreatePage{} },
		),
	)

	router.MustAddRouteExact(
		"/trainings",
		vgrouter.RouteHandlerFunc(
			func(*vgrouter.RouteMatch) { root.Body = &pages.TrainingsPage{} },
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
