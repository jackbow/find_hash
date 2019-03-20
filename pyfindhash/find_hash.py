from hashlib import md5

def main():

    # '||', 'or', 'OR', 'Or', 'oR' -- followed by a non-zero num
    vals = ["277c7c273", "276f72273", "274f52273", "274f72273", "276f52273"]

    # set to this number so it doesn't need to run for hours (should return in ~1 minute)
    pw_int = 100000000000000000000000000000538511092
    # should output 100000000000000000000000000000538611092

    found = False
    while not found:
        # convert to string, get hash
        pw_str = str(pw_int)
        h = md5(pw_str).hexdigest()
        for i in range(0, len(h)-len(vals[0]), 2):
            for v in vals:
                h_idx = i
                for j in range(len(v)):
                    if v[j] != h[h_idx]:
                        break
                    if j == len(v)-1:
                        if 0 < int(h[h_idx+1], 16) < 10:
                            print(pw_str)
                            found = True
                            break
                    h_idx += 1
                if found:
                    break
            if found:
                break
        pw_int +=1

if __name__ == "__main__":
    main()
