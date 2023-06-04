
# Tomar la ruta actual
ruta=$(pwd)

# Moverse a la ruta del repositorio Git
cd $ruta
sudo rmmod cpu_grupo18
sudo rmmod mem_grupo18
make clean
make all
sudo insmod cpu_grupo18.ko
sudo insmod mem_grupo18.ko
cat /proc/cpu_grupo18
cat /proc/mem_grupo18

# ! EJECUTAR EN LA VM GCP
# cd "/home/alvaro24_ingenieria/so2_practica1_18/module"
# bash eje.sh

# ! EJECUTAR STRESS
# stress --cpu 1 --io 1 --vm 1 --vm-bytes 128M --timeout 10s
# stress --cpu 2 --timeout 20s
# stress --vm 1 --vm-bytes 2000M --timeout 10s

