#include <iostream>
#include <string>
#include <cstring>

#include <unistd.h>
#include <fcntl.h>
#include <arpa/inet.h>

#include "udp_client.h"

UDPClient::UDPClient()
{
	memset((char*)&servaddr, 0, sizeof(servaddr));
	servaddr.sin_family = AF_INET;
}

int UDPClient::connectToServer(const char* addr, const int port)
{
	servaddr.sin_port = htons(port);
	servaddr.sin_addr.s_addr = inet_addr(addr);
	addr_size = sizeof(servaddr);
	client_sock = socket(PF_INET, SOCK_DGRAM, 0);
	if (fcntl(client_sock, F_SETFL, O_NONBLOCK) < 0) {
		return -1;
	}

	return client_sock;
}

ssize_t UDPClient::sendToServer(const std::string& s)
{
	size_t send_size = s.length() + 1;
	if (send_size < BUF_SIZE) {
		return sendto(client_sock, s.c_str(), send_size, 0, (struct sockaddr *)&servaddr, addr_size);	
	} else {
		size_t num_blocks = ((send_size % BUF_SIZE) == 0) ? (send_size / BUF_SIZE) : ((send_size / BUF_SIZE) + 1);
		for (unsigned int i = 0; i < num_blocks; ++i) {
			size_t block_size = (i == (num_blocks - 1)) ? (send_size - (num_blocks * BUF_SIZE)) : BUF_SIZE;
			ssize_t err;
			if((err = sendto(client_sock, s.c_str() + (i * BUF_SIZE), send_size, 0, (struct sockaddr *)&servaddr, addr_size)) < 0) {
				return err; 
			}
		}
	}

	return send_size;
}

ssize_t UDPClient::receiveFromServer(std::string& out)
{
	ssize_t ret;
	if((ret = recvfrom(client_sock, recv_buffer, BUF_SIZE, 0, NULL, NULL)) < 0) {
		if (errno == EWOULDBLOCK)
			return -1;
		else
			return 0;
	}
	out = std::string(recv_buffer, ret);
	return ret;
}
