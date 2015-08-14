// https://groups.google.com/forum/#!topic/golang-nuts/-N9Toa6qlu8
// shows that you can extract attribute values directly from tag/key path.
// NOTE: attribute values are encoded by prepending a hyphen, '-'.

package main

import (
	"fmt"
	"github.com/jaehue/simple-chat/Godeps/_workspace/src/github.com/clbanning/x2j"
)

var doc = `
	<doc>
		<some_tag>
			<geoInfo>
				<city name="SEATTLE"/>
				<state name="WA"/>
				<country name="USA"/>
    		</geoInfo>
			<geoInfo>
				<city name="VANCOUVER"/>
				<state name="BC"/>
				<country name="CA"/>
    		</geoInfo>
			<geoInfo>
				<city name="LONDON"/>
				<country name="UK"/>
    		</geoInfo>
		</some_tag>
	</doc>`

func main() {
	values, err := x2j.ValuesFromTagPath(doc, "doc.some_tag.geoInfo.country.-name")
	if err != nil {
		fmt.Println("err:", err.Error())
	}
	for _, v := range values {
		fmt.Println("v:", v)
	}
}
