#include <linux/module.h>
// para usar KERN_INFO
#include <linux/kernel.h>
// Header para los macros module_init y module_exit
#include <linux/init.h>
// Header necesario porque se usara proc_fs
#include <linux/proc_fs.h>

/* for copy_from_user */
#include <asm/uaccess.h>
/* Header para usar la lib seq_file y manejar el archivo en /proc*/
#include <linux/seq_file.h>
// implementacion de sysinfojiffie
// #include <linux/sysinfo.h>
// implementacion de sched para obtener el uso de CPU
#include <linux/sched.h>
#include <linux/sched/signal.h>
#include <linux/jiffies.h>
#include <linux/types.h>
#include <asm/uaccess.h>
#include <linux/mm.h>
#include <linux/time.h>
#include <linux/hugetlb.h>
#include <linux/fs.h>
#include <linux/cred.h>
#include <linux/uidgid.h>
#include <linux/delay.h>

// error:
#include <linux/cpumask.h>
#include <linux/fs.h>
#include <linux/interrupt.h>
#include <linux/kernel_stat.h>
#include <linux/sched/stat.h>
#include <linux/slab.h>
#include <linux/time_namespace.h>
#include <linux/irqnr.h>
#include <linux/sched/cputime.h>
#include <linux/tick.h>

#define PROC_NAME "cpu_202000194"

MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("obtain CPU information");
MODULE_AUTHOR("Alvaro Emmanuel Socop Perez");

static unsigned long long get_total_time(struct task_struct *task)
{
    struct task_struct *child;
    unsigned long long total_time = task->utime + task->stime;

    // Traverse the child process list recursively
    list_for_each_entry(child, &task->children, sibling)
    {
        total_time += get_total_time(child);
    }

    return total_time;
}

static int escribir_archivo(struct seq_file *archivo, void *v)
{
    struct task_struct *task;
    struct task_struct *task_hijo;
    struct list_head *children;
    long memproc;
    long memproc2;
    int indext = 0; // indice para el nombre de proceso
    struct file *file;
    struct file *file2;
    char *strstate = ""; // variable para guardar el estado del proceso
    char buffer[256];
    int len;
    long cpu_usage = 20;
    struct sysinfo info;
    long mem_usage;
    bool first = true; // solo para el primer proceso la coma
    long memoria_total = 0;
    // variables para guardar cantidad de procesos
    long int ejecucion = 0;
    long int suspendido = 0;
    long int detenido = 0;
    long int zombie = 0;
    long int totales = 0;

    //! ------------------------------- CALCULO DEL CPU -------------------------------

    unsigned long long total_time_prev = 0;
    unsigned long long used_time_prev = 0;
    unsigned long long cpu_percent = 0;
    unsigned long long jiffies_start, jiffies_end;

    jiffies_start = jiffies;
    total_time_prev = jiffies_start;

    for_each_process(task)
    {
        // total_time_prev += get_total_time(task);
        used_time_prev += task->utime + task->stime;
    }

    // Sleep for 1 second
    msleep(1000);

    jiffies_end = jiffies;

    unsigned long long total_time = 0;
    unsigned long long used_time = 0;

    // Traverse the task list to calculate total and used CPU time
    for_each_process(task)
    {
        total_time += get_total_time(task);
        used_time += task->utime + task->stime;
    }

    // Calculate the CPU percentage
    if (total_time > total_time_prev)
    {
        unsigned long long total_time_diff = total_time - total_time_prev;
        unsigned long long used_time_diff = used_time - used_time_prev;

        cpu_percent = (used_time_diff * 100) / total_time_diff;
    }

    printk(KERN_INFO "Real CPU Percent: %llu%%\n", cpu_percent);

    // printk(KERN_INFO "CPU Percent: %d%%\n", cpu_usage);

    si_meminfo(&info);

    // total_mem = (info.totalram * info.mem_unit) >> 10;  // ! memoria total en MB
    // printk(KERN_INFO "Total memory: %lu mB\n", total_mem/1000);
    memoria_total = (info.totalram * info.mem_unit) >> 10;
    // printk(KERN_INFO "Total memory: %lu MB\n", (memoria_total/1000000));
    seq_printf(archivo, "{\n");
    seq_printf(archivo, "\"cpu_usage\":"); //* "cpu_usage": 25.35,
    seq_printf(archivo, "%lu , \n", cpu_usage);
    seq_printf(archivo, "\"data\": {"); //* "data": { "proceso1":{"pid": 254, ... , "procesoshijos": [...]"}, "proceso2":{...}, ... },
    for_each_process(task)
    {
        if (!first)
        {
            seq_printf(archivo, ",");
        }
        //! 0 : ejecutando
        //! 4 : zombie
        //! 8 : detenido
        //! 1 o 1026 : suspendido
        if (task->mm)
        {
            memproc = (get_mm_rss(task->mm) << (PAGE_SHIFT - 10));
            // printk(KERN_INFO "Memoria de %s: %lu MB", task->comm, memproc);
            mem_usage = ((memproc * 100) / (memoria_total >> 10)); //! PORCENTAJE CON 2 DECIMALES PARSEAR EN FRONT
            // printk(KERN_INFO "Porcentaje de memoria de %s: %lu %%\n", task->comm,mem_usage);
        }
        if (task->state == 0 || task->state == 1026 || task->state == 2)
        {
            ejecucion++;
            strstate = "ejecucion";
        }
        else if (task->state == 4)
        {
            zombie++;
            strstate = "zombie";
        }
        else if (task->state == 8 || task->state == 8193)
        {
            detenido++;
            strstate = "detenido";
        }
        else if (task->state == 1 || task->state == 1026)
        {
            suspendido++;
            strstate = "suspendido";
        }
        totales++;
        /* Get the passwd structure for the UID */
        // char *nombre_usuario = get_cred_username(task->real_cred);

        seq_printf(archivo, "\"%d_%s\": {\"pid\": %d, \"nombre\": \"%s\", \"usuario\": \"%d\", \"estado\": \"%s\", \"ram\": %lu, \n\"procesoshijos\": [",
                   indext,
                   task->comm,
                   task->pid,
                   task->comm,
                   task->cred->uid,
                   strstate, mem_usage);
        indext++;
        task_lock(task);
        children = &(task->children);
        list_for_each_entry(task_hijo, children, sibling)
        {
            if (task_hijo->mm)
            {
                // memproc2 = (get_mm_rss(task_hijo->mm)<<PAGE_SHIFT)/(1024*1024); // ! memoria de cada proceso hijo
                // mem_usage = (memproc2*10000 / (long)(memoria_total/1000000));

                memproc = (get_mm_rss(task_hijo->mm) << (PAGE_SHIFT - 10));
                // printk(KERN_INFO "Memoria de %s: %lu MB", task->comm, memproc);
                mem_usage = ((memproc * 100) / (memoria_total >> 10));
            }
            /* Get the passwd structure for the UID */
            // pw = getpwuid(task_hijo->cred->uid.val);
            seq_printf(archivo, "{\"pid\": %d, \"nombre\": \"%s\", \"usuario\": \"%d\", \"estado\": \"%s\", \"ram\": %lu}",
                       task_hijo->pid,
                       task_hijo->comm,
                       task_hijo->real_cred->uid,
                       strstate,
                       mem_usage);

            if (task_hijo->sibling.next != &task->children)
            {
                seq_printf(archivo, ",");
            }
        }
        task_unlock(task);
        seq_printf(archivo, "]\n}");
        first = false;
    }

    seq_printf(archivo, "}, \n");
    seq_printf(archivo, "\"ejecucion\":");
    seq_printf(archivo, "%li , \n", ejecucion);

    seq_printf(archivo, "\"zombie\":");
    seq_printf(archivo, "%li , \n", zombie);

    seq_printf(archivo, "\"detenido\":");
    seq_printf(archivo, "%li , \n", detenido);

    seq_printf(archivo, "\"suspendido\":");
    seq_printf(archivo, "%li , \n", suspendido);

    seq_printf(archivo, "\"totales\":");
    seq_printf(archivo, "%li  \n", totales);
    seq_printf(archivo, "}");

    return 0;
}

// Funcion que se ejecuta cuando se le hace un cat al modulo.
static int al_abrir(struct inode *inode, struct file *file)
{
    return single_open(file, escribir_archivo, NULL);
}

// Si el su Kernel es 5.6 o mayor
static struct proc_ops operaciones =
    {
        .proc_open = al_abrir,
        .proc_read = seq_read};

static int _insert(void)
{
    proc_create("cpu_202000194", 0, NULL, &operaciones);
    printk(KERN_INFO "202000194\n");
    return 0;
}

static void _remove(void)
{
    remove_proc_entry("cpu_202000194", NULL);
    printk(KERN_INFO "Sistemas Operativos 1\n");
}

module_init(_insert);
module_exit(_remove);

// https://stackoverflow.com/questions/33594124/why-is-the-process-state-in-task-struct-stored-as-type-long