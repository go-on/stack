// Copyright (c) 2015 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package stack creates a fast and flexible middleware stack for http.Handlers.

Accepted middleware

See UseXXX methods of Stack

Batteries included

Middleware can be found in the sub package stack/middleware and stack/thirdparty.

Credits

Initial inspiration came from Christian Neukirchen's
rack for ruby some years ago.

Adapters come from carbocation/interpose (https://github.com/carbocation/interpose/blob/master/adaptors)
*/
package stack
