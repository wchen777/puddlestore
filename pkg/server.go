package pkg

import (
	"context"
	"net"
	"sync"
)

type PuddleStoreServerInstance struct {
	Mutex   sync.Mutex         // server mutex
	Addr    net.Addr           // addr of server
	Cluster *Cluster           // the puddlestore cluster
	Clients map[string]*Client // map of client ID to Client
}

func (s *PuddleStoreServerInstance) GetAddr() net.Addr {
	return s.Addr
}

func (s *PuddleStoreServerInstance) InitCluster() error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	cluster, err := CreateCluster(DefaultConfig())

	if err != nil {
		return err
	}

	s.Cluster = cluster

	return nil
}

/*
	When a client connects, we callthe cluster's create client function,
	and store the client within the server's clients map.
*/
func (s *PuddleStoreServerInstance) ClientConnect(ctx context.Context) (*ClientID, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	client, err := s.Cluster.NewClient()

	if err != nil {
		return nil, err
	}

	s.Clients[client.GetID()] = &client

	return &ClientID{
		Id: client.GetID(),
	}, nil

}

/*
	delegate to the client's exit function
*/
func (s *PuddleStoreServerInstance) ClientExit(ctx context.Context, ID *ClientID) *Success {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	client := *s.Clients[ID.Id]

	client.Exit()

	delete(s.Clients, ID.Id)

	return &Success{
		Ok: true,
	}

}

/*
	delegate to the client's open function
*/
func (s *PuddleStoreServerInstance) ClientOpen(ctx context.Context, om *OpenMessage) (*OpenResponse, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	client := *s.Clients[om.ClientId.Id]

	fd, err := client.Open(om.Filepath, om.Create, om.Write)

	return &OpenResponse{
		Success: &Success{
			Ok: err == nil, // if err is nil, then success is true
		},
		Fd: int32(fd), // return the fd
	}, err

}

/*
	deletegate to the client's close function
*/
func (s *PuddleStoreServerInstance) ClientClose(ctx context.Context, cm *CloseMessage) (*Success, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	client := *s.Clients[cm.ClientId.Id]

	err := client.Close(int(cm.Fd))

	return &Success{
		Ok: err == nil, // if err is nil, then success is true
	}, err
}

/*
	deletegate to the client's read function
*/
func (s *PuddleStoreServerInstance) ClientRead(ctx context.Context, rm *ReadMessage) (*ReadResponse, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	client := *s.Clients[rm.ClientId.Id]

	buf, err := client.Read(int(rm.Fd), uint64(rm.Offset), uint64(rm.Size))

	return &ReadResponse{
		Success: &Success{
			Ok: err == nil, // if err is nil, then success is true
		},
		Data: buf,
	}, err
}

/*
	deletegate to the client's write function
*/
func (s *PuddleStoreServerInstance) ClientWrite(ctx context.Context, wm *WriteMessage) (*Success, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	client := *s.Clients[wm.ClientId.Id]

	err := client.Write(int(wm.Fd), uint64(wm.Offset), wm.Data)

	return &Success{
		Ok: err == nil, // if err is nil, then success is true
	}, err
}

/*
	deletegate to the client's mkdir function
*/
func (s *PuddleStoreServerInstance) ClientMkdir(ctx context.Context, mdm *MkdirMessage) (*Success, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	client := *s.Clients[mdm.ClientId.Id]

	err := client.Mkdir(mdm.Path)

	return &Success{
		Ok: err == nil, // if err is nil, then success is true
	}, err
}

func (s *PuddleStoreServerInstance) ClientRemove(ctx context.Context, rmd *RemoveMessage) (*Success, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	client := *s.Clients[rmd.ClientId.Id]

	err := client.Remove(rmd.Path)

	return &Success{
		Ok: err == nil, // if err is nil, then success is true
	}, err
}

func (s *PuddleStoreServerInstance) ClientList(ctx context.Context, lmd *ListMessage) (*ListResponse, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	client := *s.Clients[lmd.ClientId.Id]

	files, err := client.List(lmd.Path)

	return &ListResponse{
		Success: &Success{
			Ok: err == nil, // if err is nil, then success is true
		},
		Result: files,
	}, err
}

/*


   rpc ClientConnect(Empty) returns (ClientID);

   rpc ClientExit(ClientID) returns (Success);

   rpc ClientOpen(OpenMessage) returns (Success);

   rpc ClientClose(CloseMessage) returns (Success);

   rpc ClientWrite(WriteMessage) returns (Success);

   rpc ClientRead(ReadMessage) returns (ReadResponse);

   rpc ClientMkdir(MkdirMessage) returns (Success);

   rpc ClientRemove(RemoveMessage) returns (Success);

   rpc ClientList(ListMessage) returns (ListResponse);

*/
