#!/bin/bash
set -e

CONTEXT=$1
ENTITY=$2
HDL=$3

BUILDDIR="../../../build/afbd/$CONTEXT/$ENTITY/c-sync/$HDL/"
IFACEDIR="../../../tests/co-simulations/common/c-sync/"
LOGDIR="/tmp/afbd/$CONTEXT/$ENTITY/c-sync/"
FIFOSPATH="/tmp/afbd/$CONTEXT/$ENTITY/"
SRCDIR="../../../tests/co-simulations/$CONTEXT/$ENTITY/c-sync/"

cp ${IFACEDIR}cosim_iface.* $BUILDDIR
cp ${SRCDIR}* $BUILDDIR

cd $BUILDDIR
gcc -Wall *.c afbd/*.c -o tb

mkdir -p $LOGDIR

./tb ${FIFOSPATH}c-sync_${HDL} \
	${FIFOSPATH}${HDL}_c-sync \
	> ${LOGDIR}${HDL}.log 2>&1 &
