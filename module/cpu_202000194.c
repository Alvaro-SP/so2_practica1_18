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

#define PROC_NAME "cpu_202000194"

MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("obtain CPU information");
MODULE_AUTHOR("Alvaro Emmanuel Socop Perez");

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
    // cpu_usage = (jiffies_to_msecs(clock_t_to_jiffies(ktime_get())) / 10); no sirvio
    int i, j;
    u64 user, nice, system, idle, iowait, irq, softirq, steal;
    u64 guest, guest_nice;
    u64 sum = 0;
    u64 sum_softirq = 0;
    unsigned int per_softirq_sums[NR_SOFTIRQS] = {0};
    struct timespec64 boottime;

    user = nice = system = idle = iowait =
        irq = softirq = steal = 0;
    guest = guest_nice = 0;
    getboottime64(&boottime);
    /* shift boot timestamp according to the timens offset */
    ktime_get_boottime();
    for_each_possible_cpu(i)
    {
        struct kernel_cpustat kcpustat;
        u64 *cpustat = kcpustat.cpustat;

        kcpustat_cpu_fetch(&kcpustat, i);

        user += cpustat[CPUTIME_USER];
        nice += cpustat[CPUTIME_NICE];
        system += cpustat[CPUTIME_SYSTEM];
        idle += get_idle_time(&kcpustat, i);
        // iowait += get_iowait_time(&kcpustat, i);
        irq += cpustat[CPUTIME_IRQ];
        softirq += cpustat[CPUTIME_SOFTIRQ];
        steal += cpustat[CPUTIME_STEAL];
        guest += cpustat[CPUTIME_GUEST];
        guest_nice += cpustat[CPUTIME_GUEST_NICE];
        sum += kstat_cpu_irqs_sum(i);
        sum += arch_irq_stat_cpu(i);
    }

    printk(KERN_INFO "CPU Percent: %d%%\n", cpu_usage);

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