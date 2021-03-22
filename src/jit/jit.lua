-- jit.lua

-- init
local strbuf = {}
local queryVariablesMap = {} -- generate from go runtime

-- build data header
strbuf[#buf+1] = "{\"data\":"

-- resolve SelectionSet


function resolveUser(user)
    strbuf[#buf+1] = "\"Id\":"
    strbuf[#buf+1] = user.Id
    strbuf[#buf+1] = ",\"Name\":"
    strbuf[#buf+1] = user.Name
    strbuf[#buf+1] = ",\"Email\":"
    strbuf[#buf+1] = user.Email
    strbuf[#buf+1] = ",\"Married\":"
    strbuf[#buf+1] = user.Married
    strbuf[#buf+1] = ",\"Height\":"
    strbuf[#buf+1] = user.Height
    strbuf[#buf+1] = ",\"Gender\":"
    strbuf[#buf+1] = user.Gender
    strbuf[#buf+1] = ",\"Friends\":"
    strbuf[#buf+1] = resolveFriends(user)
    strbuf[#buf+1] = ",\"Location\":"
    strbuf[#buf+1] = resolveLocation(user)
end


function resolveFriends() 


end


function resolveLocation()

end


-- build no error data end
strbuf[#buf+1] = ",\"errors\":null,\"jit-result\":true}"

-- finally, return result to go runtime
local r = table.concat(strbuf)