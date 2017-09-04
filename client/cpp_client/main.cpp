#include <string>
#include <iostream>

#include "udp_client.h"

int main()
{
	UDPClient client;
	std::string msg = "hello from client";
	std::string out_msg;
	client.connectToServer("192.168.1.237", 1203);
	client.sendToServer(msg);
	client.receiveFromServer(out_msg);
	std::cout << "Received: " << out_msg << std::endl;
	return EXIT_SUCCESS;
}
