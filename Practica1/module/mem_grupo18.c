
#include <linux/module.h>   // kernel modules
#include <linux/kernel.h>   // KERN_INFO
#include <linux/init.h>     //
#include <linux/seq_file.h> // var seq_file
#include <linux/proc_fs.h>  // procs file

#include <linux/mm.h> //total ram

#define PROC_NAME "mem_grupo18"

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Grupo18");
MODULE_DESCRIPTION("Memory");

struct sysinfo infsys;

// ciclo para obtener los procesos
static void getMemoria(struct seq_file *archivo)
{
    // Total de memoria
    long memorytotal;
    si_meminfo(&infsys);
    memorytotal = infsys.totalram * infsys.mem_unit; // bytes
    seq_printf(archivo, "\t\t\"MEMORIA_TOTAL\":%ld,\n", memorytotal);

    long memorylibre;
    memorylibre = infsys.freeram * 1000; // bytes
    seq_printf(archivo, "\t\t\"MEMORIA_LIBRE\":%ld,\n", memorylibre);

    long memorybuffer;
    memorybuffer = infsys.bufferram * 1000; // bytes
    seq_printf(archivo, "\t\t\"BUFFER\":%ld\n,", memorybuffer);

    seq_printf(archivo, "\t\t\"MEM_UNIT\":%d\n,", infsys.mem_unit);

    long porcentaje;
    porcentaje = (((memorytotal) - (infsys.freeram * infsys.mem_unit) - (infsys.bufferram * infsys.mem_unit) - (infsys.sharedram * infsys.mem_unit)) * 10000) / (memorytotal);
    seq_printf(archivo, "\t\t\"PORCENTAJE\":%ld\n", porcentaje);
}

// Escribiendo la info
static int write_file(struct seq_file *archivo, void *v)
{
    seq_printf(archivo, "{\n"); // INI.JSON
    getMemoria(archivo);
    seq_printf(archivo, "}\n"); // FIN.JSON
    return 0;
}

// Se realiza la escritura del archivo
static int to_open(struct inode *inode, struct file *file)
{
    return single_open(file, write_file, NULL);
}

static struct proc_ops operations =
    {
        // proc_open --: por error de distr.
        // proc_read --: por error de distr.
        .proc_open = to_open,
        .proc_read = seq_read};

static int iniciar_init(void)
{
    proc_create("mem_grupo18", 0, NULL, &operations);
    printk(KERN_INFO "Hola mundo, somos el grupo 18 y este es el monitor de memoria\n");
    return 0; // 0 = ERROR DE CARGA
}

static void finalizar_end(void)
{
    remove_proc_entry("mem_grupo18", NULL);
    printk(KERN_INFO "Sayonara mundo, somos el grupo 18 y este fue el monitor de memoria\n");
}

module_init(iniciar_init);
module_exit(finalizar_end);
