# What is it?

`msgdate` will scan a directory for message files in a format `[[:digit:]]{12}.*\.ext` and rename them according the timestamp in the `Date: ` message header and requested timezone.

# Flags

-dir=".": A directory name to scan

-ext=".eml": A file extension (including dot) to be recognised as a message file

-loc="America/Chicago": A location name from the IANA Time Zone database

# Example output

		$ msgdate
		main.go:117: 120906222541_to.eml => 120906132541_to.eml
		main.go:117: 120906222541_to.pdf => 120906132541_to.pdf
		main.go:117: 130429070242_to App.config => 130428220242_to App.config
		main.go:117: 130429070242_to TCP.pdf => 130428220242_to TCP .pdf
		main.go:117: 130429070242_to.eml => 130428220242_to.eml
		main.go:117: 130429070242_to.pdf => 130428220242_to.pdf
		main.go:117: 130501020325_fr App.exe => 130430170325_fr App.exe
		main.go:117: 130501020325_fr.eml => 130430170325_fr.eml
