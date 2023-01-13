package main

import (
	"github.com/Nesquiko/swimlogs/pkg/view/pages"
	"github.com/vugu/vgrouter"
	"github.com/vugu/vugu"
)

// OVERALL APPLICATION WIRING IN vuguSetup
func vuguSetup(buildEnv *vugu.BuildEnv, eventEnv vugu.EventEnv) vugu.Builder {

	router := vgrouter.New(eventEnv)

	// MAKE OUR WIRE FUNCTION POPULATE ANYTHING THAT WANTS A "NAVIGATOR".
	buildEnv.SetWireFunc(func(b vugu.Builder) {
		if c, ok := b.(vgrouter.NavigatorSetter); ok {
			c.NavigatorSet(router)
		}
	})

	// CREATE THE ROOT COMPONENT
	root := &pages.Root{}
	buildEnv.WireComponent(root) // WIRE IT
	router.MustAddRouteExact("/",
		vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
			root.Body = &pages.LandingPage{}
		}))

	// router.MustAddRouteExact("/cart",
	// 	vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
	// 		root.Body = &pages.Cart{}
	// 	}))
	//
	// router.MustAddRouteExact("/checkout",
	// 	vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
	// 		root.Body = &pages.Checkout{}
	// 	}))
	//
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
