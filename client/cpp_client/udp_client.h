#ifndef __UDP_CLIENT_H__
#define __UDP_CLIENT_H__

#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>

#define BUF_SIZE 1024

class UDPClient
{
	private:
		int client_sock;
		char recv_buffer[BUF_SIZE];
		struct hostent *hp;
		struct sockaddr_in servaddr;
		socklen_t addr_size;
	
	public:
		UDPClient();
		int connectToServer(const char*, const int);
		ssize_t sendToServer(const std::string&);	
		ssize_t receiveFromServer(std::string&);
		int onError();
		int updateServer();
};

#endif
