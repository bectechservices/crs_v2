import collect from "collect.js"

export default class WorkerPool {
    private readonly script: string;
    private readonly numberOfWorkers: number;
    private workers: Array<Worker>;

    constructor(script: string, numberOfWorkers: number) {
        this.script = script;
        this.numberOfWorkers = numberOfWorkers;
        this.workers = [];
    }

    static hasWebWorkerSupport(): boolean {
        return Boolean((window as any).Worker);
    }

    public startWorkers(): void {
        for (let i = 0; i < this.numberOfWorkers; i++) {
            this.workers.push(new Worker(this.script));
        }
    }

    public stopWorkers(): void {
        this.workers.forEach((worker) => worker.terminate())
    }

    private shareWorkLoadForWorkers(data: Array<any>): Array<Array<any>> {
        const collection = collect(data);
        return collection.chunk((data.length / this.numberOfWorkers) + 1).toArray()
    }

    run<T, S>(data: Array<T>, onData: (data: Array<S>) => void, onDone: () => void): void {
        const sharedData = this.shareWorkLoadForWorkers(data);
        let waitNumber = this.numberOfWorkers;
        this.workers.forEach((worker: Worker, index: number) => {
            worker.onmessage = function (e) {
                if (e.data.type === "terminate") {
                    waitNumber -= 1;
                    if (waitNumber === 0) {
                        onDone();
                    }
                } else {
                    onData(e.data.data);
                }
            };
            worker.postMessage(sharedData[index] || []);
        });
    }
}