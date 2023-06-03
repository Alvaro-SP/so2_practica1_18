
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