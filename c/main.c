#include "com.h"
#include "gps.h"

int fd;
char read_buffer[BUFFER_SIZE];
int read_buffer_size, now;

int main(int argc, char *argv[]) {
  if (argc != 2) {
    perror("Provide only 1 argument, the device!");
    return 1;
  }

  fd = open_port(argv[1]);

  if (set_com_config(fd, 115200, 8, 'N', 1) < 0) {
    perror("set_com_config");
    return 1;
  }

  while (1) {
    memset(read_buffer, 0, BUFFER_SIZE);
    read_buffer_size = read_Buffer(fd, read_buffer);
    if (read_buffer_size > 0) {
      read_GPS_Data(read_buffer);
      parse_GpsDATA();
      if (Save_Data.ParseData_Flag == 1 && Save_Data.Usefull_Flag == 1) {
        print_Save();
        Save_Data.Usefull_Flag = 0;
        Save_Data.ParseData_Flag = 0;
      } else {
        printf("Invalid GPS data\n");
      }
    }
  }
  close(fd);

  return 1;
}
