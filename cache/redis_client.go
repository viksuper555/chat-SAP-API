package cache

//
//var rAddr = "localhost:6379"
//
//var rdb = redis.NewClient(&redis.Options{
//	Addr: rAddr,
//	//Password: rPass, // no password set
//	DB: 0, // use default DB
//})
//
//func SetUser(name string, u model.User) error {
//	if err := jsonSet(name, u); err != nil {
//		return err
//	}
//	return nil
//}
//
//func GetUser(name string) (*model.User, error) {
//	var u model.User
//	if err := jsonGet(name, &u); err != nil {
//		return nil, err
//	}
//	return &u, nil
//}
//
//func jsonSet(key string, value interface{}) error {
//	p, err := json.Marshal(value)
//	if err != nil {
//		return err
//	}
//	return rdb.Set(key, p, 0).Err()
//}
//
//func jsonGet(key string, dest interface{}) error {
//	p, err := rdb.Get(key).Bytes()
//	if err != nil {
//		return err
//	}
//	return json.Unmarshal(p, dest)
//}
