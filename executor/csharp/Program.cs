// Copyright 2015 gRPC authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

using System;
using System.Threading.Tasks;
using System.Text;
using Grpc.Core;
using Pb;


namespace FxServer
{
    class FxImpl : FxWatcher.FxWatcherBase
    {
        // Server side handler of the SayHello RPC
        public override Task<Reply> Call(Request request, ServerCallContext context)
        {
	    // Handler Class Call
	    Fx.Function fx = new Fx.Function();
            object res = fx.Handler(request.Input); 
            return Task.FromResult(new Reply { Output = res.ToString() });
        }
    }
    
    class Program
    {
        const int Port = 50051;

        public static void Main(string[] args)
        {
            Server server = new Server
            {
                Services = { FxWatcher.BindService(new FxImpl()) },
                Ports = { new ServerPort("localhost", Port, ServerCredentials.Insecure) }
            };
            server.Start();
            //server.ShutdownAsync().Wait();
        }
    }
}
