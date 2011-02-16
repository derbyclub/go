#!/usr/bin/env bash
# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Generate Go code listing errors and other #defined constant
# values (ENAMETOOLONG etc.), by asking the preprocessor
# about the definitions.

unset LANG
export LC_ALL=C
export LC_CTYPE=C

GCC=gcc

uname=$(uname)

includes_Linux='
#define _LARGEFILE_SOURCE
#define _LARGEFILE64_SOURCE
#define _FILE_OFFSET_BITS 64
#define _GNU_SOURCE

#include <sys/types.h>
#include <sys/epoll.h>
#include <sys/inotify.h>
#include <sys/ioctl.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <linux/ptrace.h>
#include <linux/wait.h>
#include <linux/if_tun.h>
#include <net/if.h>
#include <netpacket/packet.h>
'

includes_Darwin='
#define _DARWIN_C_SOURCE
#define KERNEL
#define _DARWIN_USE_64_BIT_INODE
#include <sys/types.h>
#include <sys/event.h>
#include <sys/socket.h>
#include <sys/sockio.h>
#include <sys/sysctl.h>
#include <sys/wait.h>
#include <net/if.h>
#include <net/route.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <netinet/ip_mroute.h>
'

includes_FreeBSD='
#include <sys/types.h>
#include <sys/event.h>
#include <sys/socket.h>
#include <sys/sockio.h>
#include <sys/sysctl.h>
#include <sys/wait.h>
#include <net/if.h>
#include <net/route.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <netinet/ip_mroute.h>
'

includes='
#include <sys/types.h>
#include <fcntl.h>
#include <dirent.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <netinet/ip6.h>
#include <netinet/tcp.h>
#include <errno.h>
#include <sys/signal.h>
#include <signal.h>
'

ccflags=""
next=false
for i
do
	if $next; then
		ccflags="$ccflags $i"
		next=false
	elif [ "$i" = "-f" ]; then
		next=true
	fi
done

# Write godefs input.
(
	indirect="includes_$(uname)"
	echo "${!indirect} $includes"
	echo
	echo 'enum {'

	# The gcc command line prints all the #defines
	# it encounters while processing the input
	echo "${!indirect} $includes" | $GCC -x c - -E -dM $ccflags |
	awk '
		$1 != "#define" || $2 ~ /\(/ {next}

		$2 ~ /^E([ABCD]X|[BIS]P|[SD]I|S|FL)$/ {next}  # 386 registers
		$2 ~ /^(SIGEV_|SIGSTKSZ|SIGRT(MIN|MAX))/ {next}
		$2 ~ /^(SCM_SRCRT)$/ {next}
		$2 ~ /^(MAP_FAILED)$/ {next}

		$2 !~ /^ETH_/ &&
		$2 !~ /^EPROC_/ &&
		$2 !~ /^EQUIV_/ &&
		$2 !~ /^EXPR_/ &&
		$2 ~ /^E[A-Z0-9_]+$/ ||
		$2 ~ /^SIG[^_]/ ||
		$2 ~ /^IN_/ ||
		$2 ~ /^(AF|SOCK|SO|SOL|IPPROTO|IP|IPV6|TCP|EVFILT|EV|SHUT|PROT|MAP|PACKET|MSG|SCM|IFF|NET_RT|RTM|RTF|RTV|RTA|RTAX|MCL)_/ ||
		$2 == "SOMAXCONN" ||
		$2 == "NAME_MAX" ||
		$2 == "IFNAMSIZ" ||
		$2 == "CTL_NET" ||
		$2 == "CTL_MAXNAME" ||
		$2 ~ /^TUN(SET|GET|ATTACH|DETACH)/ ||
		$2 ~ /^(O|F|FD|NAME|S|PTRACE)_/ ||
		$2 ~ /^SIOC/ ||
		$2 !~ "WMESGLEN" &&
		$2 ~ /^W[A-Z0-9]+$/ {printf("\t$%s = %s,\n", $2, $2)}
		$2 ~ /^__WCOREFLAG$/ {next}
		$2 ~ /^__W[A-Z0-9]+$/ {printf("\t$%s = %s,\n", substr($2,3), $2)}

		{next}
	' | sort

	echo '};'
) >_const.c

# Pull out just the error names for later.
errors=$(
	echo '#include <errno.h>' | $GCC -x c - -E -dM $ccflags |
	awk '$1=="#define" && $2 ~ /^E[A-Z0-9_]+$/ { print $2 }' |
	sort
)

echo '// mkerrors.sh' "$@"
echo '// MACHINE GENERATED BY THE COMMAND ABOVE; DO NOT EDIT'
echo
godefs -c $GCC "$@" -gsyscall "$@" _const.c

# Run C program to print error strings.
(
	/bin/echo "
#include <stdio.h>
#include <errno.h>
#include <ctype.h>
#include <string.h>

#define nelem(x) (sizeof(x)/sizeof((x)[0]))

enum { A = 'A', Z = 'Z', a = 'a', z = 'z' }; // avoid need for single quotes below

int errors[] = {
"
	for i in $errors
	do
		/bin/echo '	'$i,
	done

	# Use /bin/echo to avoid builtin echo,
	# which interprets \n itself
	/bin/echo '
};

static int
intcmp(const void *a, const void *b)
{
	return *(int*)a - *(int*)b;
}

int
main(void)
{
	int i, j, e;
	char buf[1024];

	printf("\n\n// Error table\n");
	printf("var errors = [...]string {\n");
	qsort(errors, nelem(errors), sizeof errors[0], intcmp);
	for(i=0; i<nelem(errors); i++) {
		e = errors[i];
		if(i > 0 && errors[i-1] == e)
			continue;
		strcpy(buf, strerror(e));
		// lowercase first letter: Bad -> bad, but STREAM -> STREAM.
		if(A <= buf[0] && buf[0] <= Z && a <= buf[1] && buf[1] <= z)
			buf[0] += a - A;
		printf("\t%d: \"%s\",\n", e, buf);
	}
	printf("}\n\n");
	return 0;
}

'
) >_errors.c

$GCC $ccflags -o _errors _errors.c && $GORUN ./_errors && rm -f _errors.c _errors _const.c
