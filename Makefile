# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

NAME=pgskail
VERSION=$(shell cat VERSION)

local:
	godep go build -a -ldflags "-w -X main.Version $(VERSION)-dev" -o build/pgskail

release:
	rm -rf build release && mkdir build release
	for os in linux freebsd darwin ; do \
		GOOS=$$os ARCH=amd64 godep go build -a -ldflags "-w -X main.Version $(VERSION)" -o build/pgskail-$$os-amd64 ; \
		tar --transform 's|^build/||' --transform 's|-.*||' -czvf release/pgskail-$(VERSION)-$$os-amd64.tar.gz build/pgskail-$$os-amd64 README.md LICENSE ; \
	done
	GOOS=windows ARCH=amd64 godep go build -a -ldflags "-w -X main.Version $(VERSION)" -o build/pgskail-$(VERSION)-windows-amd64.exe
	zip release/pgskail-$(VERSION)-windows-amd64.zip build/pgskail-$(VERSION)-windows-amd64.exe README.md LICENSE && \
		echo -e "@ build/pgskail-$(VERSION)-windows-amd64.exe\n@=pgskail.exe"  | zipnote -w release/pgskail-$(VERSION)-windows-amd64.zip
	go get github.com/progrium/gh-release/...
	gh-release create pires/$(NAME) $(VERSION) \
		$(shell git rev-parse --abbrev-ref HEAD) $(VERSION)

clean:
	rm -rf build release

.PHONY: release clean