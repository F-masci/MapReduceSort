# MapReduceSort

Repository per lo sviluppo dell'esercizio di Sistemi Distribuiti e Cloud Computing (a.a. 2024/2025) presso l'università di Tor Vergata - autore: Francesco Masci

## Panoramica

Questo progetto propone una soluzione di ordinamento scalabile per grandi dataset, adottando il paradigma MapReduce e implementando il sistema in Go.
Il protocollo RPC è utilizzato per la comunicazione tra i nodi del cluster, garantendo un'efficiente elaborazione parallela dei dati.

Il sistema adotta un'architettura master-worker in cui un nodo master coordina l'esecuzione delle operazioni di map.
Al termine della fase di map, i dati intermedi vengono partizionati e successivamente inviati ai reducer per essere processati.
Una volta processati, i dati vengono salvati dai reducer su dei file distinti con il nome del client, il timestamp della richiesta e l'indice del reducer che ha eseguito il processamento.

I nodi worker, i cui indirizzi sono dinamicamente configurati in un file JSON condiviso, eseguono dunque le operazioni di mappatura e riduzione sui dati assegnati.
Questo approccio distribuito permette di suddividere il carico di lavoro e di elaborare grandi volumi di dati in parallelo.

## Componenti

- **Client**: Effettua la richiesta al master con i dati da ordinare.
- **Master**: Coordina le operazioni di map inviando i chunk da ordinare ai vari worker.
- **Worker**: Esegue le operazioni di map e reduce. Ogni mapper riceve una porzione dei dati, li ordina e li invia ai reducer per poterli unire.

## Struttura dei file

- **config/**: Contiene i file JSON di configurazione per il master, i mapper e i reducer.
  - `master.json`: Configurazione degli indirizzi dei master:
  ```json
    [
        {
            "host": "localhost",
            "port": 45978,
            "proto": "tcp"
        },
        {
            "host": "localhost",
            "port": 45979,
            "proto": "tcp"
        }
    ]
    ```
  - `mapper.json`: Configurazione degli indirizzi dei worker mapper:
  ```json
    [
        {
            "host": "localhost",
            "port": 45980,
            "proto": "tcp"
        },
        {
            "host": "localhost",
            "port": 45981,
            "proto": "tcp"
        }
    ]
    ```
  - `reducer.json`: Configurazione degli indirizzi dei worker reducer:
  ```json
    [
        {
            "host": "localhost",
            "port": 45989,
            "proto": "tcp"
        },
        {
            "host": "localhost",
            "port": 45990,
            "proto": "tcp"
        }
    ]
    ```

- **structs/**: Contiene i file di definizione per le strutture utilizzate.
    - `address.go`: Strutture per gli indirizzi dei nodi.
    - `request.go`: Strutture per le richieste RPC.
    - `response.go`: Strutture per le risposte RPC.

- **utils/**: Contiene funzioni utili per il caricamento della configurazione, gestione degli errori, e scrittura dei risultati.

## Utilizzo

### 0. Inizializzazione

Eseguire il comando
```bash
    go mod tidy
```
per scaricare le dipendenze di Go

###  1. Configurazione
Per un corretto avvio del sistema, è necessario configurare preventivamente i file JSON situati nella cartella ***config***.
In questi file dovranno essere inseriti gli indirizzi IP e le porte dei nodi master e worker che parteciperanno al sistema.

### 2. Master

Per avviare un'istanza di master utilizzando una delle configurazioni presenti nel file ***master.json*** si può usare il flag ```--idx``` e specificare l'indice della configurazione da utilizzare:

```bash
  go run master.go --idx 0
```

Se invece si vuole specificare manualmente la porta su cui eseguire il master, possiamo usare il flag ```--port```

```bash
  go run master.go --port 45978
```

In questo modo il master sarà in ascolto su ```tcp:localhost:45978```.

Se si vuole avviare il master su un indirizzo differente e/o su un protocollo differente basterà specificarlo tramite flag.
Supponendo di voler metter il master in ascolto su ```tcp:192.168.0.11:50900```:

```bash
  go run master.go --address 192.168.0.11 --port 50900 --proto tcp
```

### 3. Worker

Come per il master, per avviare un nodo worker va specificato tramite gli stessi flag l'indirizzo su cui il worker si trova in ascolto.
Inoltre, utilizzando i flag ```--map``` e ```--reduce``` si specifica se il nodo worker deve abilitare il servizio di map, di reduce o entrambi:

- **mapper** in ascolto su ```tcp:localhost:45980```:
```bash
  go run worker.go --address "localhost" --port 45980 --proto "tcp" --map
```

- **reducer** in ascolto su ```tcp:localhost:45989```:
```bash
  go run worker.go --address "localhost" --port 45989 --proto "tcp" --reduce
```

- **worker** che esegue sia **map** che **reduce** in ascolto su ```tcp:localhost:45000```:
```bash
  go run worker.go --address "localhost" --port 45000 --proto "tcp" --map --reduce
```

***Attenzione: affinché il nodo master possa comunicare con i worker, è necessario aggiornare i file JSON di configurazione con gli indirizzi dei worker avviati.***

### 4. Client

Durante l'avvio del client deve essere specificato l'identificativo del client (che verrà usato per il nome del file di output) tramite il flg ```--client``` e il master a cui collegarsi.
Il master può essere specificato tramite indice, per prendere la configurazione dal file JSON, oppure manualmente:

- **client** che si collega al **master** con indice di configurazione ```0```:
```bash
  go run client.go --client fmasci --master-idx 0
```

- **client** che si collega al **master** manualmente:
```bash
  go run client.go --client fmasci --master-address localhost --master-port 45978 --master-proto tcp
```

## Esempio

Nel file ```example.*``` si trova uno script di esempio per avviare il sistema usando la configurazione di default del repository.
Lo script avvia in background i nodi master e worker. Al termine dell'esecuzione sarà necessario avviare solamente il client.

- Unix:
```bash
  ./example.sh
```

- Windows:
```bash
  ./example.bat
```

Per terminare tutti i processi associati al sistema si può usare ```end_system.*```:

- Unix:
```bash
  ./end_system.sh
```

- Windows:
```bash
  ./end_system.bat
```
