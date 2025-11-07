package grab

/*

Grab is meant to provide a means to have someone markup a golang struct
similar as we do with json, but with "slp"

This is to make it easier to, from go, load a configuration file


type Person struct {
	Name string
	Age  int
	Company string
}

person, err := grab.Grab[Person](fileReader)



// --- person:

; Grab takes all paramters to extract for the struct and makes
; lambdas to call on them like functions from slpx code

(Name "Bob")
(Age 45)
(Company "Acme Corp, LLC")


//


type Example struct {
	X map[string]string
	Y []int
}


// setting a map and list example:

(X '(
    '("key1", "value1") .  ; Note how we also have to quote the inner lists
	'("key2", "value2"))
)

(Y '(1 2 3 4 5 6))


complex objects are quited manually to ensure that when processed
they become
*/
