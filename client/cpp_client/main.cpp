#include <string>
#include <iostream>
#include <future>
#include <atomic>
#include <chrono>

#include "udp_client.h"

static std::atomic<bool> send_ready(false);
static std::string msg = "";

int talk_to_server(UDPClient& uc)
{
	std::string out_msg = "";
	while (1) {
		if (send_ready.load()) {
			//std::cout << "msg is " << msg << std::endl;
			uc.sendToServer(msg);
			msg = "";
			send_ready.store(false);
		} else {
			//uc.sendToServer("");
		}
		
		ssize_t rc = uc.receiveFromServer(out_msg);

		if (rc > 0) {
			std::cout << "received " << out_msg << std::endl << ">> " << std::flush;
		} else if (rc == 0){
			return 0;
		}
	}
}


int main(int argc, char **argv)
{
	if (argc != 2) {
		std::cout << "USAGE: ./client <ip_address>" << std::endl;
		return EXIT_SUCCESS;
	}

	UDPClient client;
	client.connectToServer(argv[1], 1203);
	
	std::future<int> msg_future = std::async(std::launch::async, talk_to_server, std::ref(client));

	std::cout << ">> ";

	for (;;) {
		//std::cout << ">> ";
		std::getline(std::cin, msg);
		std::this_thread::sleep_for(std::chrono::milliseconds(1));
		send_ready.store(true);
		std::this_thread::sleep_for(std::chrono::milliseconds(1));
	}

	int result = msg_future.get();

	return EXIT_SUCCESS;
}
