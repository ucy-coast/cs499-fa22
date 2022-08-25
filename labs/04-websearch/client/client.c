#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <sys/socket.h>
#include <resolv.h>
#include <netinet/tcp.h>
#include <netinet/in.h>
#include <netdb.h>
#include <pthread.h>
#include <stdlib.h>
#include <string>
#include <iostream>
#include <fstream>
#include <sstream>
#include <sys/time.h>
#include <ctime>
#include <signal.h>
#include <vector>
#include <unistd.h>

#define MAX_QUERY_SIZE 100
#define MAX_REQUEST_SIZE 256
#define RECEIVE_BUFFER_SIZE 1500
#define DUMMY_BIG_VALUE 1000000
#define CHUNK_SIZE_NOT_FOUND -1

using namespace std;

typedef struct
{
        int client_num;
        int start;
        int finish;
//      int part_size;
} start_finish;

typedef unsigned long long uint64;

pthread_mutex_t START_LOCK;
//ofstream log;
//pthread mutex_t LOG_LOCK;
string * QUERIES_TO_EXECUTE; //global query to execute array
double * RESPONSE_TIMES;  //global response time array
//long * HITS_ARRAY; //for saving hits return per search. Validation that service is the same across setups
vector<string> responseVector ; //DEBUG ONLY
vector <struct sockaddr_in> serv_addresses (1);
struct hostent * server;
int TCP_NO_DELAY_FLAG;
string SERVER_IP;
string REQUEST_TYPE;
int PORTNO;
int KEEP_ALIVE;



void error(const char *msg);

double getTime();

void  * send_requests(start_finish * bounds);

int getChunkSize(string & line);

void error(const char *msg)
{
    perror(msg);
   // exit(0);
}

double getTime() //copied from http://stackoverflow.com/questions/1861294/how-to-calculate-execution-time-of-a-code-snippet-in-c
{
 struct timeval tv;

 gettimeofday(&tv, NULL);

  double ret = tv.tv_usec;
 /* Convert from micro seconds (10^-6) to milliseconds (10^-3) */
 ret /= 1000;

 /* Adds the seconds (10^0) after converting them to milliseconds (10^-3) */
 ret += (tv.tv_sec * 1000);

 return ret;
}
extern int ONEPORT; //TODO this is buggy doing just to remember what extern means hahahahaahahha
void *  send_requests(start_finish * bounds){
        pthread_mutex_lock(&START_LOCK);
        //fprintf(stderr,"Client number %d started with start %d and end %d\n",bounds->client_num,bounds->start,bounds->finish);
        int sockfd = socket(AF_INET, SOCK_STREAM, 0);
        if (sockfd < 0)
        error("ERROR opening socket");

         int result = setsockopt(sockfd,            /* socket affected */
                                 IPPROTO_TCP,     /* set option at TCP level */
                                 TCP_NODELAY,     /* name of option */
                                 (char *) &TCP_NO_DELAY_FLAG,  /* the cast is historical
                                                         cruft */
                                 sizeof(int));    /* length of option value */
         if (result < 0)
        {
                perror("unable to set TCP_NODELAY\n");
                exit(1);
        }

        if(REQUEST_TYPE=="onlyHitsPrintIds.jsp")
        {
                struct timeval timeout;
                timeout.tv_sec = 120;
                timeout.tv_usec = 0;


                if (setsockopt (sockfd, SOL_SOCKET, SO_RCVTIMEO, (char *)&timeout,
                        sizeof(timeout)) < 0)
                        error("setsockopt failed\n");

                /*if (setsockopt (sockfd, SOL_SOCKET, SO_SNDTIMEO, (char *)&timeout,
                        sizeof(timeout)) < 0)
                error("setsockopt failed\n");*/
        }

        int addr_index=0;
        if(ONEPORT==0)
        {
                addr_index=bounds->client_num;
        }

        if (connect(sockfd,(struct sockaddr *) &serv_addresses.at(addr_index/*bounds->client_num*/),sizeof(serv_addresses.at(addr_index))) < 0) //Todo at i when multiple tomcat instances watchout
        error("ERROR connecting");
        //fprintf(stderr,"Client number %d succesfully connected\n",bounds->client_num);
        pthread_mutex_unlock(&START_LOCK);

        //Now start sending queries
        string send_buffer ;
        string keep_alive;
        char receive_buffer [RECEIVE_BUFFER_SIZE];
        int headers =0; //only for debug
        int requests =0;
        int fail = 0;
        bool failed = false;
        if(KEEP_ALIVE==1)
        {
                keep_alive="Connection: keep-alive\r\n";
        }
        else
        {
                keep_alive="Connection: close\r\n"; //TODO when not using keep alive something is messed up. In general sometimes that stats are not shown wtf?
        }



        for(int i=bounds->start;i<bounds->finish; i++)
        {
                ostringstream stringStream;
                stringStream << "GET /"<<REQUEST_TYPE<<"?query="<<QUERIES_TO_EXECUTE[i]<<" HTTP/1.1\r\nHost: "<<SERVER_IP<<":"<<(PORTNO+addr_index)<<"\r\n"<<keep_alive/*<< "Accept-Charset: US-ASCII\r\n"*/ << "\r\n";
                string send_buffer = stringStream.str();
                requests++;
                //cerr<<"send "<<requests<<"\n"; //DEBUG
                //cerr<<send_buffer; //Only for validation
                //cerr<<"client "<< bounds->client_num << " send query " << i <<"\n"; //only for validation
                /*pthread_mutex_lock(&START_LOCK);
                if(i==bounds->start)
                {
                        //log<<i+bounds->client_num*bounds->part_size<<"\n";
                        log<<i<<"\n";
                }
                else
                {
                        //log<<(i+bounds->client_num*bounds->part_size)-1<<"\n"<<i+bounds->client_num*bounds->part_size<<"\n";
                        log<<i-1<<"\n"<<i<<"\n";
                }
                pthread_mutex_unlock(&START_LOCK);*/
                double startTime = getTime();//start measuring time
                int n = write(sockfd,send_buffer.c_str(),send_buffer.length());
                if (n < 0)
                {
                        pthread_mutex_lock(&START_LOCK);
                        error("ERROR writing to socket"); //When setting maxKeepAlive in tomcat server xml to -1 the probelms seem to be solved. But still I need to understatn how to re-initiate connection with tomcat and what to do when keep-alive is no used
                        close(sockfd);
                        if (sockfd < 0)
                        {
                         error("ERROR opening socket");
                        }


                        if (connect(sockfd,(struct sockaddr *) &serv_addresses.at(addr_index),sizeof(serv_addresses.at(addr_index))) < 0) //Todo at i when multiple tomcat instances watchout
                        {
                                error("ERROR connecting");
                        }
                        pthread_mutex_unlock(&START_LOCK);
                }

                bzero(receive_buffer,RECEIVE_BUFFER_SIZE);
                bool readingHeader = true; //Alway the first msg part should be header
                vector <string> header ;
                bool endOfHtml=false;
                bool readingChunk = false;
                int chunkSize=-1;
                int totalChunkRead=0;
                string chunkSizeLine="";
                readAgain:
                while( (n = read(sockfd,receive_buffer,RECEIVE_BUFFER_SIZE-1 )) >0) //TODO decide if you need the -1
                {
                        //cerr<<"DEBUG Read "<<n<<"bytes from server\n"; //DEBUG
                        //responseVector.push_back(receive_buffer);
                        if(REQUEST_TYPE!="search.jsp") //TODO make this the correct way
                        {
                                //cerr <<"#"<< i << " " << receive_buffer; //DEBUG
                                responseVector.at(i) = receive_buffer;//DEBUG
                                  //responseVector.push_back(receive_buffer);
                                /*char * pch = strtok (receive_buffer,"#");
                                int token = 0;
                                while (pch != NULL && token < 1)
                                {
                                    pch = strtok (NULL,"#");
                                    token++;
                                }
                                HITS_ARRAY[i] = atol(pch); */
                                break;
                        }

                        istringstream inputStream (receive_buffer);
                        string line;

                        if(readingHeader)
                        {
                                if(header.empty()) //Check if header is empty.  In case the header comes partialy we don't want to count it twice.
                                {
                                        headers++;
                                }


                        //      cerr<<"\nHeader "<<headers<<" start\n\n"; //DEBUG
                                while(!inputStream.eof()) //read header line by line
                                {
                                        getline(inputStream,line);
                                        bool endLine = true; //I mean header end line
                                        for(int i=0;i<line.length();i++)
                                        {
                                                if(!isspace(line[i]))
                                                {
                                                        header.push_back(line);
                                                        //cerr<<line; //Cerr line header for debug
                                                        endLine=false;
                                                        break;
                                                }
                                        }
                                        if(endLine)
                                        {
                                                //cerr<<"This is the breaking line:"<<line<<"and consists of ";
                                                /*for(int i=0;i<line.length();i++)
                                                {
                                                        cerr<<(int)line[i]<<" "<<"\n";
                                                }*/
                                                readingHeader=false;
                                                /*while(!inputStream.eof())
                                                {
                                                        getline(inputStream,line);
                                                        cerr<<line;
                                                }*/
                                                //We have all the header lets process it
                                                if(!header.empty())
                                                {
                                                        if(header.at(0).find("HTTP/1.1 200 OK")==string::npos)
                                                        {
                                                                failed=true;
                                                                fail++;
                                                                cerr<<"Fail requests with response "<<header.at(0)<<"header number " << headers ;

                                                        }
                                                /*for(int i=0;i< header.size(); i++) shoud be for recoginziing chunked encoding content size blaha blah blah
                                                {
                                                        if(header.at(i).find("Transfer-Encoding: chunked")!=string::npos )
                                                        {
                                                                cerr << "\nchecking for find correctness\n";
                                                                cerr<<header.at(i);
                                                        }
                                                }*/
                                                        header.clear();
                                                }
                                                break;
                                        }
                                }
                                //cerr<<"\nHeader "<<headers<<" end\n\n";
                        }


                        /*if(strstr(receive_buffer,"Server: Apache-Coyote/1.1") != NULL) //Header found ONLY for debug at least for now.
                        {

                        }*/
                        bool chunkedResponse=true; //TODO fix this. This should be read from header
                        if(!readingHeader)
                        {
                                if(chunkedResponse)
                                {

                                        while(!inputStream.eof())
                                        {

                                                if(readingChunk)
                                                {
                                                        char c;
                                                        //cerr<<"READING BODY OF CHUNK"<<endl; //DEBUG
                                                        while(totalChunkRead<chunkSize)
                                                        {       c =inputStream.get();
                                                                //cerr<<c; //DEBUG
                                                                totalChunkRead++;
                                                        }
                                                        //cerr<<"Reading body before CR"<<endl; //DEBUG
                                                        while((c=inputStream.get())!='\r' && !inputStream.eof() ) //eat /r/n
                                                        {
                                                        //      cerr<<c; //DEBUG
                                                        }
                                                        //cerr<<"End of Reading body before CR"<<endl; //DEBUG
                                                        if(c!='\r' && inputStream.eof()) //This means that I didn't managed to read the whole chunk so I need to complete it at the next iteration
                                                        {
                                                        //      cerr<<"PARTIALY READED BODY CHUNK\n"; //DEBUG
                                                                continue;
                                                        }
                                                        //cerr<<"Reading body before LF"<<endl; //DEBUG
                                                        while((c=inputStream.get())!='\n' && !inputStream.eof() ) //eat /r/n/. TODO Lets hope the two of them don't come partially. For now it seems okay but generally should fix this
                                                        {
                                                        //              cerr<<c; //DEBUG
                                                        }
                                                        //cerr<<"END of Reading body before LF"<<endl; //DEBUG
                                                        //cerr<<"END OF CHUNK BODY"<<endl; //DEBUG
                                                        //if(totalChunkRead>=chunkSize)
                                                        //{
                                                                //cerr<<"finished reading chunk size "<<chunkSize<<" with total bytes read "<<totalChunkRead<<endl; //DEBUG
                                                        readingChunk=false;
                                                        //}
                                                }
                                                else
                                                {
                                                        char c;
                                                        //cerr<<"trying to read chunk size"<<endl; //DEBUG
                                                        while( (c=inputStream.get())!='\n' && !inputStream.eof() ) //eat until /r
                                                        {
                                                        //      cerr<<(int)c<<" "; //DEBUG
                                                                chunkSizeLine.push_back(c);
                                                        //getline(inputStream,line);
                                                        }
                                                        if(c!='\n' && inputStream.eof()) //This means that I didn't managed to read a whole line I need to complete it at the next iteration
                                                        {
                                                        //      cerr<<"PARTIALY READED CHUNK SIZE\n"; //DEBUG
                                                                continue;
                                                        }


                                                        //cerr<<"end of reading chunk size"<<endl; //DEBUG
                                                        //cerr<<"chunkLineSize is "<<chunkSizeLine; //DEBUG
                                                                chunkSize = getChunkSize(chunkSizeLine);
                                                                //cerr<<"new chunk size "<<chunkSize<<endl; //DEBUG
                                                                if(chunkSize == 0)
                                                                {
                                                                //      cerr<<"line "<<line<<" ended response\n"; //DEBUG
                                                                        endOfHtml=true;
                                                                        break;
                                                                }
                                                                else if(chunkSize!=CHUNK_SIZE_NOT_FOUND)
                                                                {
                                                                        totalChunkRead=0;
                                                                        readingChunk=true;
                                                                }


                                                                 chunkSizeLine.clear(); //Clear the string for accepting the next line

                                                }
                                        }
                                }
                        }
                        if(endOfHtml)
                        {
                         break;
                        }
                        bzero(receive_buffer,RECEIVE_BUFFER_SIZE);
                } //End of while read loop

                if (n <= 0)
                {
                        pthread_mutex_lock(&START_LOCK);
                        if (n==0)
                        {
                                cerr<<"Read 0 bytes from server\n";
                                //goto readAgain;

                        }
                        else
                        {
                                error("ERROR reading from socket");
                        }
                        pthread_mutex_unlock(&START_LOCK);
                        if(!endOfHtml)
                        {
                        //      goto readAgain;
                                bzero(receive_buffer,RECEIVE_BUFFER_SIZE);
                        }
                }

                //if(!failed) //todo
                //{
                        RESPONSE_TIMES[i] = getTime() - startTime; //record the response time in ms
                        //cout<<RESPONSE_TIMES[i];
                //}
        }
        return NULL;
}



int getChunkSize(string & line) // -1 no otherwise chunk size in bytes//TODO maybe need to check also for chunk extensions
{
        if(!isdigit(line[0]) && line[0]!='a' && line[0]!='b' && line[0]!='c' && line[0]!='d' && line[0]!='e' && line[0]!='f')
        {
                return -1;
        }
        //bool number=true;
        int i;
        string number="";
        number.push_back(line[0]);
        for( i=1; i< line.size(); i++)
        {
                if(!isdigit(line[i]) && line[i]!='a' && line[i]!='b' && line[i]!='c' && line[i]!='d' && line[i]!='e' && line[i]!='f')
                {
                        break;
                }
                number.push_back(line[i]);

        }
        i++;
        for(; i< line.size(); i++)
        {
                if(!isspace(line[i]))
                {
                        return -1;
                }
        }
//      cerr<<"chunk size hex is "<<number<<endl; //DEBUG
        int x;
        std::stringstream ss;
        ss << std::hex << number;
        ss >> x;
        //cerr<<"chunk size dec is " << x<<endl; //DEBUG
/*      if(x==0)
        {
                cerr<<"The 0 line contains";
                for(int i=0; i< number.size(); i++)
                {
                        cerr<<(int)number[i];
                }
                cerr<<"END of line 0\n";
        }*///DEBUG
        return x;
}

/*void cntl_c_handler(int sig)
{
  pthread_t  f1_thread;
  int i1=1;
  pthread_create(&f1_thread,NULL,(void *)&f1,&i1);
}*/

void quickSort(int left, int right);
int partition(int left, int right);
int queriesToSort;
//quick sort


int partition(int left, int right)
{
    double pivot_element = RESPONSE_TIMES[left];
    int lb = left, ub = right;
    double temp;

    while (left < right)
    {
        while(RESPONSE_TIMES[left] <= pivot_element && left < queriesToSort)
            left++;
        while(RESPONSE_TIMES[right] > pivot_element  && right > 0)
            right--;
        if (left < right)
        {
            temp        = RESPONSE_TIMES[left];
            RESPONSE_TIMES[left]  = RESPONSE_TIMES[right];
            RESPONSE_TIMES[right] = temp;
        }
    }
    RESPONSE_TIMES[lb] = RESPONSE_TIMES[right];
    RESPONSE_TIMES[right] = pivot_element;
    return right;
}



void quickSort(int left, int right)
{
    if (left < right)
    {
        int pivot = partition(left, right);
        quickSort(left, pivot-1);
        quickSort(pivot+1, right);
    }
}

int ONEPORT=1;

//A multithreaded c client that send requests as fast as possible

int main(int argc, char * argv[]) //usage  server_ip, port, queriesFIle, queries to execute, number of client threads,request_type (search.jsp etc), TCP_NO_DELAY ON/OFF (1/0), KEEP_ALIVE (1/0) full path to responses file
{
        signal(SIGPIPE,SIG_IGN); //DON't exit on broken pipe error
        bool fast=false;
        if(argc==12)
        {
                if(strcmp(argv[11],"FAST")==0)
                {
                        fast=true;
                }
        }
        else if(argc!=11)
        {
                printf("usage  server_ip, port, queriesFIle, queries to execute, number of client threads,request_type (search.jsp etc), TCP_NO_DELAY ON/OFF (1/0), KEEP_ALIVE (1/0), full path to responses file, for multiple Frontends enter 0 otherwise 1\n");
                return 0;
        }
        if (pthread_mutex_init(&START_LOCK, NULL) != 0)
        {
                printf("\n mutex init failed\n");
                return 1;
        }
        //log.open ("log", ofstream::out);
        cerr<<"Client started\n";
        time_t start,end;
        time (&start);
        //Parameter reading
        SERVER_IP = argv[1];
        PORTNO = atoi(argv[2]);
        ifstream query_file (argv[3]);
        if (!query_file) {
                cout <<argv[3] << " File does not exists\n";
                return 1;
        }
        int num_queries_to_execute = atoi(argv[4]);
        int clients = atoi(argv[5]);
        REQUEST_TYPE = argv[6];
        TCP_NO_DELAY_FLAG = (atoi(argv[7]));
        KEEP_ALIVE = atoi(argv[8]);
        if(TCP_NO_DELAY_FLAG !=1)
        {
                TCP_NO_DELAY_FLAG=0;
        }


        //Boring client-server stuff
         server = gethostbyname(SERVER_IP.c_str());
        if (server == NULL) {
        fprintf(stderr,"ERROR, no such host\n");
        exit(0);
        }

        serv_addresses.resize(clients);

        ONEPORT=atoi(argv[10]);

        if(ONEPORT==0)
        {
                for(int i=0;i<serv_addresses.size();i++)
                {
                bzero((char *) &serv_addresses.at(i), sizeof(serv_addresses.at(i)));
                serv_addresses.at(i).sin_family = AF_INET;
                bcopy((char *)server->h_addr,
                (char *)&serv_addresses.at(i).sin_addr.s_addr,
                 server->h_length);
                serv_addresses.at(i).sin_port = htons(PORTNO+i);
                }
                //PORTNO=atoi(argv[2]);
        }
        else
        {
                //for(int i=0;i<serv_addresses.size();i++)
                //{
                int i=0;
                        bzero((char *) &serv_addresses.at(i), sizeof(serv_addresses.at(i)));
                        serv_addresses.at(i).sin_family = AF_INET;
                        bcopy((char *)server->h_addr,
                        (char *)&serv_addresses.at(i).sin_addr.s_addr,
                        server->h_length);
                        serv_addresses.at(i).sin_port = htons(PORTNO);
                //}
        }

        /*
                struct sockaddr_in  serv_addr;
                bzero((char *) &serv_addr, sizeof(serv_addr));
        serv_addr.sin_family = AF_INET;
        bcopy((char *)server->h_addr,
         (char *)&serv_addr.sin_addr.s_addr,
         server->h_length);
        serv_addr.sin_port = htons(PORTNO);
        */

        QUERIES_TO_EXECUTE   = new string [num_queries_to_execute];
        int i=0;
        int real_num_of_queries = 0; //In case given input with actual file size does not match
        for(i=0; i<num_queries_to_execute; i++) //Read the queries file
        {
                 getline (query_file, QUERIES_TO_EXECUTE[i]);
                 if( query_file.eof() ) //while the end of file is NOT reached
                 {
                                break;
                 }
                 real_num_of_queries++;
        }
        query_file.close();
        time (&end);
        cerr<<"Reading queries from file took " << difftime(end,start) <<" sec\n";
        /*test that query read correctly (It okay)

        for(i=0; i<real_num_of_queries; i++)
        {
                cout << queries_to_execute[i] <<"\n";
        }*/



        RESPONSE_TIMES = (double *) malloc(sizeof(double)*real_num_of_queries); //create the response times array
        for(i=0;i<real_num_of_queries;i++)
        {
                RESPONSE_TIMES[i]=DUMMY_BIG_VALUE;
        }

        /*if(REQUEST_TYPE=="onlyHits.jsp") //TODO make this the correct way
        {
                HITS_ARRAY = (long *) malloc(sizeof(long)*real_num_of_queries); //create the hits array
        }*/

        pthread_t * clients_pid = (pthread_t *) malloc(sizeof(pthread_t)*clients); //create the thread tids array
        start_finish * bounds  = (start_finish *) malloc (sizeof(start_finish)*clients);
        int part = real_num_of_queries / clients;
        int remainder = real_num_of_queries % clients;

        responseVector.resize(real_num_of_queries);

        //start measuring the throughput
        cerr<<"Starting the run\n";

        time (&start);
        for(i=0; i<clients-1; i++)
        {
                bounds[i].client_num = i;
                bounds[i].start = i * part;
                bounds[i].finish = bounds[i].start +  part;
                //bounds[i].part_size=part;
                pthread_create(&clients_pid[i],NULL, reinterpret_cast<void* (*)(void*)>(send_requests),&bounds[i]);
        }
        bounds[i].client_num = i;
        bounds[i].start = i * part;
        bounds[i].finish = bounds[i].start +  part + remainder;
        //bounds[i].part_size=part+remainder;
        pthread_create(&clients_pid[i],NULL, reinterpret_cast<void* (*)(void*)>(send_requests),&bounds[i]);

        for(i=0; i<clients; i++)
        {
                pthread_join(clients_pid[i],NULL);
        }

        time (&end);
        double dif = difftime (end,start);
        cerr<<"Run Finish and took "<< dif << " sec  Now start processing stats\n";

        ofstream responseTimeFile (argv[9]);
        for(int m=0;m<real_num_of_queries;m++)
        {
                responseTimeFile << RESPONSE_TIMES[m]<<"\n";
        }
        responseTimeFile.close();
        if(fast)
        {
                return 0;
        }
        time(&start);

        /*bool swapped =true; //Bubble sort is really slow so I turned to quick sort
        for (int c = 0 ; swapped &&  c < ( real_num_of_queries - 1 ); c++) //first sort response times array
        {
                swapped = false;
                for (int d = 0 ; d < real_num_of_queries - c - 1; d++)
                {
                        if (RESPONSE_TIMES[d] > RESPONSE_TIMES[d+1])
                        {
                                uint64 swap       = RESPONSE_TIMES[d];
                                RESPONSE_TIMES[d]   = RESPONSE_TIMES[d+1];
                                RESPONSE_TIMES[d+1] = swap;
                                swapped=true;
                        }
                }
      }*/
        queriesToSort=real_num_of_queries;
        quickSort(0, real_num_of_queries -1);
        time(&end);
        cerr<<"Sorting response times finish" <<" and took " << difftime(end,start) <<" sec \n";

        //find where dummy values begin for not succeded quereis
        for (int c = 0 ; c <  real_num_of_queries ; c++)
        {
                if(RESPONSE_TIMES[c]==DUMMY_BIG_VALUE)
                {
                        real_num_of_queries=c; //change this variable to the actual sucedded queries and break from loop
                        break;
                }
        }

        cout<< "ops/sec " << real_num_of_queries/dif <<"\n" ;
        double average=0;

        for (int c = 0 ; c <  real_num_of_queries ; c++)
        {
                //cout<<"number "<<RESPONSE_TIMES[c] <<"\n";
                average += RESPONSE_TIMES[c];
        }

        cout<<"average " << (average/real_num_of_queries)<<"\n";  //This is one way to calculate average

        int rank50 = int(0.5 * real_num_of_queries + 0.5)-1;
        int rank90  = int(0.9 * real_num_of_queries + 0.5)-1;
        int rank99  = int(0.99 * real_num_of_queries + 0.5)-1;
        /*cout <<"rank 90 is " <<rank90<<"\n";
        cout <<"rank 99 is " <<rank99<<"\n";    */
        cout<<"50th " << RESPONSE_TIMES[rank50]<<"\n";
        cout<<"90th "  <<RESPONSE_TIMES[rank90]<<"\n";
        cout<<"99th "  <<RESPONSE_TIMES[rank99]<<"\n";
        cout<<"succeded_queries "<<real_num_of_queries<<"\n";

        if(REQUEST_TYPE=="onlyHits.jsp" || REQUEST_TYPE=="onlyHitsRex.jsp" || REQUEST_TYPE=="onlyHitspana.jsp") //TODO make this the correct way
        {
                cout<<"printing hits\n";
                double average = 0;
                for (int i=0; i<responseVector.size(); i++)
                {
                        //cout<<HITS_ARRAY[i]<<"\n";
                        cerr<<responseVector.at(i); //DEBUG
                        //average += HITS_ARRAY[i];
                }
                cout<<"end of printing hits. average number of hits is "<<(average/real_num_of_queries)<<"\n";
        }

        cerr.flush();
        cout.flush();
        cerr<<"Bye\n";
        pthread_mutex_destroy(&START_LOCK);
        //log.close();
        return 0;
}




//Left garbish code
//FILE * query_file = fopen(argv[3],"r" );
//      char **  queries_to_execute  = (char **) malloc(sizeof(char *) * num_queries_to_execute);
        //remove the trailing \n
                /*int j=0;
                for(j=0;j<MAX_QUERY_SIZE; j++)
                {
                        if(queries_to_execute[i][j] == '\n')
                        {
                                queries_to_execute[i][j]='\0';
                                break;
                        }
                }*/

                        /*queries_to_execute[i] = (char*) malloc(sizeof(char)* MAX_QUERY_SIZE);
                bzero(queries_to_execute[i],MAX_QUERY_SIZE);
                if(fgets(queries_to_execute[i],MAX_QUERY_SIZE,query_file) ==NULL) //Some error happened or EOF encountered
                {
                        break;
                }*/

//sprintf(send_buffer,"GET /%s?query=%s HTTP/1.1\n\rHost: %s:%d\n\r%s\n\r",REQUEST_TYPE, QUERIES_TO_EXECUTE[i] ,SERVER_IP,PORTNO,keep_alive);
