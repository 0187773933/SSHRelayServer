package secretbox

import (
	"fmt"
	"time"
	"os"
	"strings"
	"encoding/base64"
	"io/ioutil"
	math_rand "math/rand"
	crypto_rand "crypto/rand"
	// bcrypt "golang.org/x/crypto/bcrypt"
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
func generate_random_bytes( length int ) ( result []byte ) {
	result = make( []byte , length )
	_ , err := crypto_rand.Read( result )
	if err != nil { panic( err ) }
	number_of_extra_shuffles := math_rand.Intn( shuffle_max - shuffle_min ) + shuffle_min
	for i := 0; i < number_of_extra_shuffles; i++ {
		result = shuffle_bytes( result )
	}
	return
}

type Box struct {
	Key []byte `json:"key"`
	Nonce []byte `json:"nonce"`
	HumanString string `json:"human_string"`
}
func New() ( box Box ) {
	box = Box{}
	box.Key = generate_random_bytes( key_size )
	box.Nonce = generate_random_bytes( nonce_size )
	human_string_bytes := []byte( fmt.Sprintf( "%s@@@===@@@%s" , string( box.Key ) , string( box.Nonce ) ) )
	box.HumanString = base64.StdEncoding.EncodeToString( human_string_bytes )
    // bcrypted_human_string , _ := bcrypt.GenerateFromPassword( human_string_bytes , bcrypt.DefaultCost )
    // fmt.Println( string( bcrypted_human_string ) )
	return
}
func Load( secretbox_key_nonce_b64 string ) ( box Box ) {
	secretbox_key_nonce_bytes , secretbox_key_nonce_bytes_err := base64.StdEncoding.DecodeString( secretbox_key_nonce_b64 )
	if secretbox_key_nonce_bytes_err != nil { panic( secretbox_key_nonce_bytes_err ) }
	parts := strings.Split( string( secretbox_key_nonce_bytes ) , "@@@===@@@" )
	box = Box{}
	box.Key = []byte( parts[ 0 ] )
	box.Nonce = []byte( parts[ 1 ] )
	box.HumanString = base64.StdEncoding.EncodeToString( []byte( fmt.Sprintf( "%s@@@===@@@%s" , string( box.Key ) , string( box.Nonce ) ) ) )
	return
}
func ( box *Box ) get_enforced_24_byte_nonce() ( result [24]byte ) {
	copy( result[:], box.Nonce[:] )
	return
}
func ( box *Box ) get_enforced_32_byte_key() ( result [32]byte ) {
	copy( result[:], box.Key[:] )
	return
}
func ( box *Box ) SealMessage( message string ) ( result string ) {
	//t, err := base64.URLEncoding.DecodeString(cookie.Value)
	var enforced_32_byte_key [32]byte
	copy( enforced_32_byte_key[:], box.Key[:] )
	var enforced_24_byte_nonce [24]byte
	copy( enforced_24_byte_nonce[:], box.Nonce[:] )
	sealed := make( []uint8 , nonce_size )
	copy( sealed , box.Nonce[:] )
	sealed = secretbox.Seal( sealed , []byte( message ) , &enforced_24_byte_nonce , &enforced_32_byte_key );
	result = base64.StdEncoding.EncodeToString( sealed )
	return
}
func ( box *Box ) OpenMessage( message_b64 string ) ( result string ) {
	result = "failed"
	var enforced_32_byte_key [32]byte
	copy( enforced_32_byte_key[:], box.Key[:] )
	var enforced_24_byte_nonce [24]byte
	copy( enforced_24_byte_nonce[:], box.Nonce[:] )
	message_bytes , message_bytes_err := base64.StdEncoding.DecodeString( message_b64 )
	if message_bytes_err != nil { panic( message_bytes_err ) }
	decoded_bytes , ok := secretbox.Open( nil , message_bytes[nonce_size:] , &enforced_24_byte_nonce , &enforced_32_byte_key )
	if ok != true { panic( ok ) }
	result = string( decoded_bytes )
	return
}
func ( box *Box ) OpenFile( file_path string ) ( result string ) {
	result = "failed"
	data , _ := ioutil.ReadFile( file_path )
	lines := strings.Split( string( data ) , "\n" )
	lines_cleaned := []string{}
	for i := range lines { if lines[ i ] != "" { lines_cleaned = append( lines_cleaned , lines[ i ] ) } }
	file_b64_data := strings.Join( lines_cleaned , "\n" )
	nonce := box.get_enforced_24_byte_nonce()
	key := box.get_enforced_32_byte_key()
	file_bytes , file_bytes_err := base64.StdEncoding.DecodeString( file_b64_data )
	if file_bytes_err != nil { panic( file_bytes_err ) }
	decoded_bytes , ok := secretbox.Open( nil , file_bytes[nonce_size:] , &nonce , &key )
	if ok != true { panic( ok ) }
	result = string( decoded_bytes )
	return
}
func ( box *Box ) SealPlainTextToFile( file_path string , message string ) ( result string ) {
	result = "failed"
	nonce := box.get_enforced_24_byte_nonce()
	key := box.get_enforced_32_byte_key()
	sealed := make( []uint8 , nonce_size )
	copy( sealed , box.Nonce[:] )
	sealed = secretbox.Seal( sealed , []byte( message ) , &nonce , &key );
	sealed_b64_string := base64.StdEncoding.EncodeToString( sealed )
	out_file , _ := os.Create( file_path )
	defer out_file.Close()
	out_file.WriteString( sealed_b64_string )
	result = "success"
	return
}
func ( box *Box ) SealPlainFileToEncryptedFile( plain_file_path string , encrypted_file_path string ) ( result string ) {
	result = "failed"
	nonce := box.get_enforced_24_byte_nonce()
	key := box.get_enforced_32_byte_key()
	sealed := make( []uint8 , nonce_size )
	copy( sealed , box.Nonce[:] )
	data , _ := ioutil.ReadFile( plain_file_path )
	lines := strings.Split( string( data ) , "\n" )
	lines_cleaned := []string{}
	for i := range lines { if lines[ i ] != "" { lines_cleaned = append( lines_cleaned , lines[ i ] ) } }
	cleaned_plaintext := strings.Join( lines_cleaned , "\n" )
	sealed = secretbox.Seal( sealed , []byte( cleaned_plaintext ) , &nonce , &key );
	sealed_b64_string := base64.StdEncoding.EncodeToString( sealed )
	out_file , _ := os.Create( encrypted_file_path )
	defer out_file.Close()
	out_file.WriteString( sealed_b64_string )
	result = "success"
	return
}