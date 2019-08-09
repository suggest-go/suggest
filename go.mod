module github.com/alldroll/suggest

go 1.12

require (
	github.com/RoaringBitmap/roaring v0.4.19-0.20190803203039-d0ce1763c352
	github.com/alldroll/cdb v1.0.1
	github.com/alldroll/go-datastructures v0.0.0-20190322060030-1d3a19ff3b29
	github.com/edsrzf/mmap-go v0.0.0-20190108065903-904c4ced31cd
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.1
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3 // indirect
	golang.org/x/sys v0.0.0-20190426135247-a129542de9ae // indirect
)

replace github.com/RoaringBitmap/roaring => github.com/alldroll/roaring v0.4.19-0.20190803203039-b3f4f3210e08
