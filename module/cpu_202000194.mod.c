#include <linux/module.h>
#define INCLUDE_VERMAGIC
#include <linux/build-salt.h>
#include <linux/elfnote-lto.h>
#include <linux/export-internal.h>
#include <linux/vermagic.h>
#include <linux/compiler.h>

BUILD_SALT;
BUILD_LTO_INFO;

MODULE_INFO(vermagic, VERMAGIC_STRING);
MODULE_INFO(name, KBUILD_MODNAME);

__visible struct module __this_module
__section(".gnu.linkonce.this_module") = {
	.name = KBUILD_MODNAME,
	.init = init_module,
#ifdef CONFIG_MODULE_UNLOAD
	.exit = cleanup_module,
#endif
	.arch = MODULE_ARCH_INIT,
};

#ifdef CONFIG_RETPOLINE
MODULE_INFO(retpoline, "Y");
#endif


static const struct modversion_info ____versions[]
__used __section("__versions") = {
	{ 0xbdfb6dbb, "__fentry__" },
	{ 0x5436cc1d, "proc_create" },
	{ 0x92997ed8, "_printk" },
	{ 0x5b8239ca, "__x86_return_thunk" },
	{ 0x471a2c4d, "single_open" },
	{ 0x90074bc8, "filp_open" },
	{ 0x51988f6a, "kernel_read" },
	{ 0x28462739, "filp_close" },
	{ 0xbcab6ee6, "sscanf" },
	{ 0x40c7247c, "si_meminfo" },
	{ 0x6ca0c764, "seq_printf" },
	{ 0x7ba7391a, "init_task" },
	{ 0xba8fbd64, "_raw_spin_lock" },
	{ 0xb5b54b34, "_raw_spin_unlock" },
	{ 0xd0da656b, "__stack_chk_fail" },
	{ 0xbd651dae, "remove_proc_entry" },
	{ 0x4bfec637, "seq_read" },
	{ 0xb83992f2, "module_layout" },
};

MODULE_INFO(depends, "");

