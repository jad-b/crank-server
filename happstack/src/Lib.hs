{-# LANGUAGE
    OverloadedStrings
  , ScopedTypeVariables
#-}
module Lib (
    someFunc
) where
import Control.Applicative ((<$>), optional)
import Data.Maybe (fromMaybe)
import Data.Text (Text)
import Data.Text.Lazy (unpack)
import Happstack.Lite


someFunc :: IO ()
someFunc = putStrLn "someFunc"
