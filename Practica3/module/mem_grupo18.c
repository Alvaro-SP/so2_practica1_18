
#include <linux/module.h>   // kernel modules
#include <linux/kernel.h>   // KERN_INFO
#include <linux/init.h>     //
#include <linux/seq_file.h> // var seq_file
#include <linux/proc_fs.h>  // procs file

#include <linux/mm.h> //total ram

#define PROC_NAME "mem_grupo18"

// Informacion de la memoria
MODULE_LICENSE("GPL");
MODULE_AUTHOR("Grupo18");
MODULE_DESCRIPTION("Memory");

// Struct de sysinfo, donde se encuentra la información de la memoria ram
struct sysinfo infsys;

/**
    Funcion para obtener los valores de la memoria ram, que se obtienen del archivo sysinfo

    @param: seq_file archivo a leer

    @return: --
*/
static void getMemoria(struct seq_file *archivo)
{
    // Total de memoria
    long memorytotal;
    si_meminfo(&infsys);
    memorytotal = infsys.totalram * infsys.mem_unit; // bytes
    seq_printf(archivo, "\t\t\"MEMORIA_TOTAL\":%ld,\n", memorytotal);
    // Memoria libre
    long memorylibre;
    memorylibre = infsys.freeram * 1000; // bytes
    seq_printf(archivo, "\t\t\"MEMORIA_LIBRE\":%ld,\n", memorylibre);
    // Memoria compartida
    long memorybuffer;
    memorybuffer = infsys.bufferram * 1000; // bytes
    seq_printf(archivo, "\t\t\"BUFFER\":%ld\n,", memorybuffer);
    // Memoria compartida
    seq_printf(archivo, "\t\t\"MEM_UNIT\":%d\n,", infsys.mem_unit);
    // Memoria compartida
    long porcentaje;
    porcentaje = (((memorytotal) - (infsys.freeram * infsys.mem_unit) - (infsys.bufferram * infsys.mem_unit) - (infsys.sharedram * infsys.mem_unit)) * 10000) / (memorytotal);
    seq_printf(archivo, "\t\t\"PORCENTAJE\":%ld\n", porcentaje);
}

/**
    Funcion para escribir  los valores obtenidos de la memoria ram

    @param: seq_file archivo a leer

    @return: --
*/
static int write_file(struct seq_file *archivo, void *v)
{
    // se escribe el archivo JSON con la informacion de la memoria
    seq_printf(archivo, "{\n"); // INI.JSON
    getMemoria(archivo);
    seq_printf(archivo, "}\n"); // FIN.JSON
    return 0;
}

/**
 * @brief
 *
 * @param inode : estructura que contiene información sobre un archivo en el disco
 * @param file : estructura que contiene información sobre un archivo abierto
 * @return int
 */
static int to_open(struct inode *inode, struct file *file)
{
    // se abre el archivo
    return single_open(file, write_file, NULL);
}

// static struct proc_ops operations =
//     {
//         // proc_open --: por error de distr.
//         // proc_read --: por error de distr.
//         .proc_open = to_open,
//         .proc_read = seq_read};

/**
 * @brief
 *
 */
// Si el kernel es menor al 5.6 usan file_operations
static struct file_operations operaciones =
    {
        // estas funciones son las que se ejecutan al abrir el archivo
        .open = to_open,
        .read = seq_read};

static int iniciar_init(void)
{
    // se crea el archivo
    proc_create("mem_grupo18", 0, NULL, &operaciones);
    // mensaje de inicio
    printk(KERN_INFO "Hola mundo, somos el grupo 18 y este es el monitor de memoria\n");
    return 0; // 0 = ERROR DE CARGA
}

static void finalizar_end(void)
{
    // se elimina el archivo
    remove_proc_entry("mem_grupo18", NULL);
    // mensaje de despedida
    printk(KERN_INFO "Sayonara mundo, somos el grupo 18 y este fue el monitor de memoria\n");
}

// se registran las funciones de inicio y fin
module_init(iniciar_init);
module_exit(finalizar_end);
