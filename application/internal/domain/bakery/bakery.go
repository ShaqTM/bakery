package bakery

type Bakery struct{
	Storage *ports.Storage
	Log *logrus.Logger
}

func New(log *logrus.Logger, storage *ports.Storage)*Bakery{
	return Bakery{
		Log:log,
		Storage: storage
	}		


}