# Copyright 2015 Muir Manders.  All rights reserved.
# Use of this source code is governed by a MIT-style
# license that can be found in the LICENSE file.

task :default => :install

GOBIN = "#{ENV['GOPATH']}/bin"

file "#{GOBIN}/godep" do
  sh "go get github.com/tools/godep"
end

file "#{GOBIN}/go-bindata" do
  sh "go get github.com/jteeuwen/go-bindata/..."
end

task :install => ["#{GOBIN}/godep", "#{GOBIN}/go-bindata"] do
  # generate static http resources
  sh "PATH='#{GOBIN}:#{ENV['PATH']}' godep go generate ./..."
  sh "#{GOBIN}/godep go install ./..."
end

task :dev_install => ["#{GOBIN}/godep", "#{GOBIN}/go-bindata"] do
  sh "#{GOBIN}/go-bindata --debug --pkg http -o http/resources.go --prefix http/resources http/resources"
  sh "#{GOBIN}/godep go install ./..."
end
