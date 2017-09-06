#include <string>
#include <iostream>

#include "udp_client.h"

int main(int argc, char **argv)
{
	if (argc != 2) {
		std::cout << "USAGE: ./client <ip_address>" << std::endl;
		return EXIT_SUCCESS;
	}

	UDPClient client;
	std::string msg = "hello from client";
	std::string out_msg;
	client.connectToServer(argv[1], 1203);
	client.sendToServer(msg);
	client.receiveFromServer(out_msg);
	std::cout << "Received: " << out_msg << std::endl;
	return EXIT_SUCCESS;
}
