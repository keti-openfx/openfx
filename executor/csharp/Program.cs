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
using System.IO;
using System.Threading.Tasks;
using System.Text;
using Grpc.Core;
using Google.Protobuf;
using Pb;


namespace FxWatcherServer
{
    class FxWatcherImpl : FxWatcher.FxWatcherBase
    {
        // Server side handler of the SayHello RPC
        public override Task<Reply> Call(Request request, ServerCallContext context)
        {
	    // Handler Class Call
	        Fx.Function fx = new Fx.Function();
		    byte[] res = fx.Handler(request.Input.ToByteArray());
			string result = Encoding.Default.GetString(res);
            return Task.FromResult(new Reply { Output = result });
        }
    }
   
    class Program
    {
        const int Port = 50051;
        
        private static async Task Process()
        {

            while (true)
            {
                await Task.Delay(100);
            }

        }
    
        public static void Main(string[] args)
        {
            
			
			File.Create("/tmp/.lock");
			Console.WriteLine("Writing lock-file to: /tmp/.lock");

			Server server = new Server
            {
                Services = { FxWatcher.BindService(new FxWatcherImpl()) },
                Ports = { new ServerPort("0.0.0.0", Port, ServerCredentials.Insecure) }
            };
            server.Start();

            Console.WriteLine("[fxwatcher] start service.");
	    
            var processTask = Process();
            processTask.Wait();   
            server.ShutdownAsync().Wait();
        }
    }
}
