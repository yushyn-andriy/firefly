
fn pow (base, exp) {
    result = {"res": base};
    for (i = 0; i<exp; i = i + 1;) {
        r  = result["res"] * base;
        result["res"] = r;
        print("result:",r);
    }
    return result;
}

print(pow(2, 4));

