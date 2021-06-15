package secretbox

import (
    "fmt"
    "time"
    "encoding/base64"
    math_rand "math/rand"
    crypto_rand "crypto/rand"
    secretbox "golang.org/x/crypto/nacl/secretbox"
)
const key_size = 32
const nonce_size = 24
const shuffle_min = 50
const shuffle_max = 100

func shuffle_bytes( bytes []byte ) ( result []byte ) {
    result = make( []byte , len( bytes ) )
    //copy( bytes[:] , shuffled_bytes[:] )
    random_source := math_rand.New( math_rand.NewSource( time.Now().Unix() ) )
    random_indexes := random_source.Perm( len( bytes ) )
    for index , item := range random_indexes {
        result[ index ] = bytes[ item ]
    }
    return
}
func GenerateRandomBytes( length int ) ( result []byte ) {
    result = make( []byte , length )
    _ , err := crypto_rand.Read( result )
    if err != nil { panic( err ) }
    number_of_extra_shuffles := math_rand.Intn( shuffle_max - shuffle_min ) + shuffle_min
    for i := 0; i < number_of_extra_shuffles; i++ {
        result = shuffle_bytes( result )
    }
    return
}
func SealMessage( message string , nonce []byte , key []byte ) ( result string ) {
    //t, err := base64.URLEncoding.DecodeString(cookie.Value)
    var enforced_32_byte_key [32]byte
    copy( enforced_32_byte_key[:], key[:] )
    var enforced_24_byte_nonce [24]byte
    copy( enforced_24_byte_nonce[:], nonce[:] )
    sealed := make( []uint8 , nonce_size )
    copy( sealed , nonce[:] )
    sealed = secretbox.Seal( sealed , []byte( message ) , &enforced_24_byte_nonce , &enforced_32_byte_key );
    result = base64.StdEncoding.EncodeToString( sealed )
    return
}
func OpenMessage( message_b64 string , nonce []byte , key []byte  ) ( result string ) {
    result = "failed"
    var enforced_32_byte_key [32]byte
    copy( enforced_32_byte_key[:], key[:] )
    var enforced_24_byte_nonce [24]byte
    copy( enforced_24_byte_nonce[:], nonce[:] )
    message_bytes , message_bytes_err := base64.StdEncoding.DecodeString( message_b64 )
	if message_bytes_err != nil { panic( message_bytes_err ) }
	decoded_bytes , ok := secretbox.Open( nil , message_bytes[nonce_size:] , &enforced_24_byte_nonce , &enforced_32_byte_key )
	if ok != true { panic( ok ) }
    result = string( decoded_bytes )
    return
}
func main() {
    message := "testing wadu wadu wadu"
    fmt.Println( message )
    key := GenerateRandomBytes( key_size )
    //fmt.Printf( "key === %v\n" , key )
    nonce := GenerateRandomBytes( nonce_size )
    //fmt.Printf( "nonce === %v\n" , nonce )
    sealed := SealMessage( message , nonce , key )
    fmt.Println( sealed )
    opened := OpenMessage( sealed , nonce , key )
    fmt.Println( opened )
}
