# apid-core

apid-core is a library that provides a container for publishing APIs that provides core services to its plugins 
including configuration, API publishing, data access, and a local pub/sub event system.

Disambiguation: You might be looking for the executable builder, [apid](https://github.com/30x/apid).  

## Services

apid provides the following services:

* apid.API()
* apid.Config()
* apid.Data()
* apid.Events()
* apid.Log()
 
### Initialization of services and plugins

A driver process must initialize apid and its plugins like this:

    apid.Initialize(factory.DefaultServicesFactory()) // when done, all services are available
    apid.InitializePlugins() // when done, all plugins are running
    api := apid.API() // access the API service
    err := api.Listen() // start the listener


Once apid.Initialize() has been called, all services are accessible via the apid package functions as details above. 

## Plugins

The only requirement of an apid plugin is to register itself upon init(). However, generally plugins will access
the Log service and some kind of driver (via API or Events), so it's common practice to see something like this:
 
    var log apid.LogService
     
    func init() {
      apid.RegisterPlugin(initPlugin)
    }
    
    func initPlugin(services apid.Services) error {
    
      log = services.Log().ForModule("myPluginName") // note: could also access via `apid.Log().ForModule()`
      
      services.API().HandleFunc("/verifyAPIKey", handleRequest)
    }
    
    func handleRequest(w http.ResponseWriter, r *http.Request) {
      // respond to request
    }

## Running Tests

    go test $(glide novendor)
