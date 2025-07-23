import multiprocessing

def burn_cpu():
    while True:
        pass  # Infinite loop to keep the CPU core busy

if __name__ == "__main__":
    num_cores = multiprocessing.cpu_count()
    processes = []

    # for _ in range(num_cores):
    for _ in range(2):
        p = multiprocessing.Process(target=burn_cpu)
        p.start()
        processes.append(p)

    for p in processes:
        p.join()