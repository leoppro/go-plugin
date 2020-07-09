function max(num1, num2)

    if (num1 > num2) then
       result = num1;
    else
       result = num2;
    end
 
    return result;
 end


 function emit_row(row)
 
    return row.table;
 end