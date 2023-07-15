/**
 * Please implement a minimal clone of the ls program. We have chosen this
 *  exercise as it will require you to use structs, pointers and arrays, as well
 *  as some C standard library functions with interesting interfaces. It will
 *  also likely to be substantial enough to merit some degree of code
 *  organization.
 *
 * Minimally, it should list the contents of a directory including some
 *  information about each file, such as file size. As a stretch goal, use
 *  `man ls` to identify any interesting flags you may wish to support, and
 *  implement them.
 */
#include <dirent.h>
#include <fcntl.h>
#include <getopt.h>
#include <grp.h>
#include <pwd.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <time.h>

#define true 1
#define false 0

typedef int bool;
typedef void (*dirent_printer)(const struct dirent *);
typedef bool (*visibility_decider)(const char *);
void simple_printer(const struct dirent *);
void long_printer(const struct dirent *);
bool print_all_dirs(const char *);
bool omit_dotted_dirs(const char *);

struct options {
  dirent_printer printer;                 /* -l */
  visibility_decider should_print_dirent; /* -a */
};
static const char *recognized_options = "al";
static struct options opts = {/* default behavior */
                              simple_printer, omit_dotted_dirs};
struct options *init_opts(int argc, char **argv);

/**
 * USAGE: ls [-al] [dir]
 */
int main(int argc, char **argv) {
  char *dirname;
  /**
   * collect options.  The program assumes they're provided before the
   *  arguments.
   */
  if (init_opts(argc, argv) != NULL) {
    argc -= optind;
    argv += optind;
  } else {
    printf("ls: could not parse options.");
    exit(EXIT_FAILURE);
  }

  /**
   * choose directory.
   */
  if (argc == 0)
    dirname = ".";
  else
    dirname = *argv;

  DIR *dir_p = opendir(dirname);
  if (dir_p == NULL) {
    printf("ls: could not open directory %s\n.", dirname);
    exit(EXIT_FAILURE);
  }

  struct dirent *directory_entry;
  while ((directory_entry = readdir(dir_p)) != NULL) {
    if (opts.should_print_dirent(directory_entry->d_name))
      opts.printer(directory_entry);
  }

  exit(EXIT_SUCCESS);
}

void long_printer(const struct dirent *dir_entry) {
  struct stat dirent_stats;
  stat(dir_entry->d_name, &dirent_stats);

  char last_modified[BUFSIZ];
  struct tm *last_modified_tm = localtime(&dirent_stats.st_mtimespec.tv_sec);
  strftime(last_modified, BUFSIZ, "%b %d %H:%M", last_modified_tm);

  char *owner = getpwuid(dirent_stats.st_uid)->pw_name;
  char *group = getgrgid(dirent_stats.st_gid)->gr_name;

  /**
   * This could be extracted
   */
  char filetype;
  switch (dirent_stats.st_mode & S_IFMT) {
  case S_IFIFO:
    filetype = 'p';
    break;
  case S_IFCHR:
    filetype = 'c';
    break;
  case S_IFDIR:
    filetype = 'd';
    break;
  case S_IFBLK:
    filetype = 'b';
    break;
  case S_IFREG:
    filetype = '-';
    break;
  case S_IFLNK:
    filetype = 'l';
    break;
  case S_IFSOCK:
    filetype = 's';
    break;
  case S_IFWHT:
    filetype = 'w';
    break;
  default:
    printf("ls: unrecognized file mode %d", dirent_stats.st_mode);
    exit(EXIT_FAILURE);
  }

  char mode[11] = {filetype,
                   (dirent_stats.st_mode & S_IRUSR) ? 'r' : '-',
                   (dirent_stats.st_mode & S_IWUSR) ? 'w' : '-',
                   (dirent_stats.st_mode & S_IXUSR) ? 'x' : '-',
                   (dirent_stats.st_mode & S_IRGRP) ? 'r' : '-',
                   (dirent_stats.st_mode & S_IWGRP) ? 'w' : '-',
                   (dirent_stats.st_mode & S_IXGRP) ? 'x' : '-',
                   (dirent_stats.st_mode & S_IROTH) ? 'r' : '-',
                   (dirent_stats.st_mode & S_IWOTH) ? 'w' : '-',
                   (dirent_stats.st_mode & S_IXOTH) ? 'x' : '-',
                   '\0'};

  /**
   * improvement: formatting (padding etc.)
   * This would require adjusting function signature, because the number is
   * based on all entries, not just one.
   */
  printf("%s\t%2.d\t%s\t%s\t%8.lld\t%s\t%s\n", mode, dirent_stats.st_nlink,
         owner, group, dirent_stats.st_size, last_modified, dir_entry->d_name);
}
void simple_printer(const struct dirent *dir) { printf("%s\n", dir->d_name); }

bool print_all_dirs(const char *dirname) { return true; }
bool omit_dotted_dirs(const char *dirname) {
  return *dirname == '.' ? false : true;
}
struct options *init_opts(int argc, char **argv) {
  int opt;

  while ((opt = getopt(argc, argv, recognized_options)) != -1)
    switch (opt) {
    case 'l':
      opts.printer = long_printer;
      break;
    case 'a':
      opts.should_print_dirent = print_all_dirs;
      break;
    case '?':
    default:
      printf("ls: Unrecognized option.  Usage: `ls [-al] [file..]`");
      return NULL;
    }

  return &opts;
}
