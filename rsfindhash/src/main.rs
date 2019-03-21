use termion;
use num_cpus;
use md5::compute;
use rand::{thread_rng, seq::SliceRandom};
use std::{thread, process, sync::mpsc, time::Instant};

const QUOTE_BYTE: u8 = 0x27;
const MID_VALS: [[u8; 2]; 5] = [[0x7c, 0x7c],
                                [0x6f, 0x72],
                                [0x6f, 0x52],
                                [0x4f, 0x52],
                                [0x4f, 0x72]];
const ONE_BYTE: u8 = 0x31;
const NINE_BYTE: u8 = 0x39;
const PW_LEN: u8 = 5;  // length of passwords generated (94^5 = 7,339,040,224 possible strings)
const CHARSET: &[u8] =  b"ABCDEFGHIJKLMNOPQRSTUVWXYZ\
                         abcdefghijklmnopqrstuvwxyz\
                         0123456789)(*&^%$#@!~{}[]|\
                         \\\'\";:_-+=<>,./?`"; // character set to use to generate passwords
const CPU_FRAC_DENOM: usize = 1; // denominator of fraction of cpu cores to run
const PRINT_FREQ: u32 = 230000; //` 1 in x pws are printed as output

fn check_hash(hash: md5::Digest) -> bool {
    for i in 0..(hash.0.len()-5) {
        if hash[i] != QUOTE_BYTE || hash[i+3] != QUOTE_BYTE ||
           hash[i+4] < ONE_BYTE || hash[i+4] > NINE_BYTE {
               continue;
        }
        let mut or_byte = false;
        for v in &MID_VALS {
            if hash[i+1] == v[0] && hash[i+2] == v[1] {
                or_byte = true;
            }
        }
        if or_byte {
            return true;
        }
    }
    return false;
}

// https://rust-lang-nursery.github.io/rust-cookbook/algorithms/randomness.html#create-random-passwords-from-a-set-of-user-defined-characters
fn gen_string(len: u8) -> String {
    let mut rng = thread_rng();
    let string: Option<String> = (0..len)
        .map(|_| Some(*CHARSET.choose(&mut rng)? as char))
        .collect();
    return string.unwrap();
}

struct Password {
    string: String,
    n_checked: u32,
    valid: bool,
}

fn printer(rx: mpsc::Receiver<Password>) {
    let now = Instant::now();
    let nsearchers = (num_cpus::get()/CPU_FRAC_DENOM) as u32;
    print!("{}", termion::cursor::Hide);
    for msg in rx {
        println!("{}\n{} million hashes in {}sec",
            msg.string,
            msg.n_checked * nsearchers / 1000000,
            now.elapsed().as_secs());
        if msg.valid {
            println!("{}", termion::cursor::Show); // give user back their cursor
            process::exit(0);
        }
        else { // reset cursor to overwrite output if not valid pw
            print!("{}{}",
            termion::cursor::Up(2),
            termion::cursor::Left(4));
        }
    }
}

fn search(tx: mpsc::Sender<Password>) {
    let mut i: u32 = 0;
    // used to benchmark
    let nsearchers = num_cpus::get()/CPU_FRAC_DENOM;
    loop {
        i += 1;
        let pw = gen_string(PW_LEN);
        let valid = check_hash(compute(&pw));
        if valid || i % PRINT_FREQ == 0 {
            tx.send(Password{
                        string: pw.clone(),
                        n_checked: i,
                        valid: valid,
                }).unwrap();
        }
        // used to benchmark
        if i == 10000000/(nsearchers as u32) { process::exit(0); }
    }
}

fn main() {
    // quick test
    assert!(check_hash(compute("QS+02")));
    assert!(!check_hash(compute("invalid")));

    // spawn printer
    let (tx, rx) = mpsc::channel();
    let printer_thr = thread::spawn(|| printer(rx));

    // spawn searchers
    let mut searcher_thrs = vec![];
    for _ in 0..(num_cpus::get()/CPU_FRAC_DENOM) {
        let txc = tx.clone();
        searcher_thrs.push(thread::spawn(|| search(txc)));
    }

    // wait for printer
    // (searchers search endlessly)
    let _ = printer_thr.join();
}
