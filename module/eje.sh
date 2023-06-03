
# Tomar la ruta actual
ruta=$(pwd)

# Moverse a la ruta del repositorio Git
cd $ruta
sudo rmmod cpu_202000194
sudo rmmod ram_202000194
make clean
make all
sudo insmod cpu_202000194.ko
sudo insmod ram_202000194.ko
cat /proc/cpu_202000194
cat /proc/ram_202000194

# ! EJECUTAR EN LA VM GCP
# cd "/home/alvaro24_ingenieria/so2_practica1_18/module"
# bash eje.sh

# ! EJECUTAR STRESS
# stress --cpu 1 --io 1 --vm 1 --vm-bytes 128M --timeout 10s
# stress --cpu 2 --timeout 20s

