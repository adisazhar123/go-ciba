package go_ciba

//func TestRead(t *testing.T) {
//	file, err := os.Open("data/key.pem")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer file.Close()
//
//	b, err := ioutil.ReadAll(file)
//
//	block, _ := pem.Decode(b)
//	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
//	fmt.Println(key.N)
//}
//
//func TestGenerate(t *testing.T) {
//	reader := rand.Reader
//	bitSize := 2048
//
//	privateKey, _ := rsa.GenerateKey(reader, bitSize)
//	// Validate Private Key
//	err := privateKey.Validate()
//	if err != nil {
//		log.Println(err)
//	}
//
//	privateKeyFile, _ := os.Create("test_data/key.pem")
//
//	_ = pem.Encode(privateKeyFile, &pem.Block{
//		Type:  "RSA PRIVATE KEY",
//		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
//	})
//
//	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
//
//	publicKeyFile, _ := os.Create("test_data/public.pem")
//
//	_ = pem.Encode(publicKeyFile, &pem.Block{
//		Type:  "PUBLIC KEY",
//		Bytes: asn1Bytes,
//	})
//
//	defer privateKeyFile.Close()
//	defer publicKeyFile.Close()
//}
//
//func TestGenerate2(t *testing.T) {
//	Priv, _ := rsa.GenerateKey(rand.Reader, 4096)
//
//	privBytes := pem.EncodeToMemory(&pem.Block{
//		Type:  "RSA PRIVATE KEY",
//		Bytes: x509.MarshalPKCS1PrivateKey(Priv),
//	})
//
//	pubASN1, err := x509.MarshalPKIXPublicKey(&Priv.PublicKey)
//	if err != nil {
//		// do something about it
//	}
//
//	pubBytes := pem.EncodeToMemory(&pem.Block{
//		Type:  "PUBLIC KEY",
//		Bytes: pubASN1,
//	})
//
//	ioutil.WriteFile("test_data/key.pem", privBytes, 0644)
//	ioutil.WriteFile("test_data/public.pem", pubBytes, 0644)
//
//}
