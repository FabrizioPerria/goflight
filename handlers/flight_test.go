package handlers

// const (
// 	uri    = "mongodb://localhost:27017"
// 	dbName = "goflight_test"
// )

// type testdb struct {
// 	UserStore db.UserStorer
// 	Client    *mongo.Client
// }
//
// func setup() (*testdb, error) {
// 	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
// 	if err != nil {
// 		return nil, err
// 	}
// 	db := db.NewMongoDbUserStore(client, dbName)
// 	return &testdb{UserStore: db, Client: client}, nil
// }
//
// func teardown(t *testing.T, db *testdb) {
// 	if err := db.UserStore.Drop(context.Background()); err != nil {
// 		t.Fatal(err)
// 	}
// 	if err := db.Client.Disconnect(context.TODO()); err != nil {
// 		t.Fatal(err)
// 	}
// }
